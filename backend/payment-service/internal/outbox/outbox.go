package outbox

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/ErenKarakus1/Ride-Sharing-System/backend/payment-service/internal/events"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventStore struct {
	db *pgxpool.Pool
}

type Record struct {
	ID    string
	Event events.Event
}

func NewEventStore(db *pgxpool.Pool) *EventStore {
	return &EventStore{db: db}
}

func (s *EventStore) Publish(ctx context.Context, event events.Event) error {
	payload, err := json.Marshal(event.Data)
	if err != nil {
		return err
	}

	const query = `
		INSERT INTO payment_outbox_events (id, event_type, event_key, payload)
		VALUES ($1, $2, $3, $4)
	`
	_, err = s.db.Exec(ctx, query, uuid.NewString(), event.Type, event.Key, payload)
	return err
}

func (s *EventStore) Pending(ctx context.Context, limit int) ([]Record, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	const query = `
		WITH next_events AS (
			SELECT id
			FROM payment_outbox_events
			WHERE status IN ('pending', 'failed') AND attempts < 5
			ORDER BY created_at
			LIMIT $1
			FOR UPDATE SKIP LOCKED
		)
		UPDATE payment_outbox_events e
		SET attempts = attempts + 1, updated_at = now()
		FROM next_events
		WHERE e.id = next_events.id
		RETURNING e.id, e.event_type, e.event_key, e.payload
	`

	rows, err := tx.Query(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records := make([]Record, 0)
	for rows.Next() {
		var record Record
		var payload json.RawMessage
		if err := rows.Scan(&record.ID, &record.Event.Type, &record.Event.Key, &payload); err != nil {
			return nil, err
		}
		record.Event.Data = payload
		records = append(records, record)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return records, tx.Commit(ctx)
}

func (s *EventStore) MarkSent(ctx context.Context, id string) error {
	const query = `UPDATE payment_outbox_events SET status = 'sent', published_at = now(), updated_at = now() WHERE id = $1`
	_, err := s.db.Exec(ctx, query, id)
	return err
}

func (s *EventStore) MarkFailed(ctx context.Context, id string) error {
	const query = `
		UPDATE payment_outbox_events
		SET status = 'failed', updated_at = now()
		WHERE id = $1
	`
	_, err := s.db.Exec(ctx, query, id)
	return err
}

type Dispatcher struct {
	store     *EventStore
	publisher events.Publisher
}

func NewDispatcher(store *EventStore, publisher events.Publisher) *Dispatcher {
	return &Dispatcher{store: store, publisher: publisher}
}

func (d *Dispatcher) Run(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			d.flush(ctx)
		}
	}
}

func (d *Dispatcher) flush(ctx context.Context) {
	records, err := d.store.Pending(ctx, 25)
	if err != nil {
		log.Printf("failed to load payment outbox events: %v", err)
		return
	}

	for _, record := range records {
		if err := d.publisher.Publish(ctx, record.Event); err != nil {
			log.Printf("failed to publish payment outbox event id=%s type=%s: %v", record.ID, record.Event.Type, err)
			_ = d.store.MarkFailed(ctx, record.ID)
			continue
		}
		_ = d.store.MarkSent(ctx, record.ID)
	}
}
