package payment

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, payment Payment) (Payment, error) {
	const query = `
		INSERT INTO payments (id, ride_id, rider_id, driver_id, amount, currency, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, ride_id, rider_id, driver_id, amount, currency, status, created_at, updated_at
	`

	if payment.ID == "" {
		payment.ID = uuid.NewString()
	}

	return scanPayment(r.db.QueryRow(
		ctx,
		query,
		payment.ID,
		strings.TrimSpace(payment.RideID),
		strings.TrimSpace(payment.RiderID),
		payment.DriverID,
		payment.Amount,
		strings.TrimSpace(payment.Currency),
		payment.Status,
	))
}

func (r *PostgresRepository) Get(ctx context.Context, id string) (Payment, error) {
	const query = `
		SELECT id, ride_id, rider_id, driver_id, amount, currency, status, created_at, updated_at
		FROM payments
		WHERE id = $1
	`

	payment, err := scanPayment(r.db.QueryRow(ctx, query, strings.TrimSpace(id)))
	if err == pgx.ErrNoRows {
		return Payment{}, ErrPaymentNotFound
	}
	if err != nil {
		return Payment{}, err
	}

	return payment, nil
}

func (r *PostgresRepository) UpdateStatus(ctx context.Context, id string, status Status) (Payment, error) {
	const query = `
		UPDATE payments
		SET status = $2, updated_at = now()
		WHERE id = $1
		RETURNING id, ride_id, rider_id, driver_id, amount, currency, status, created_at, updated_at
	`

	payment, err := scanPayment(r.db.QueryRow(ctx, query, strings.TrimSpace(id), status))
	if err == pgx.ErrNoRows {
		return Payment{}, ErrPaymentNotFound
	}
	if err != nil {
		return Payment{}, err
	}

	return payment, nil
}

type paymentRow interface {
	Scan(dest ...any) error
}

func scanPayment(row paymentRow) (Payment, error) {
	var payment Payment
	err := row.Scan(
		&payment.ID,
		&payment.RideID,
		&payment.RiderID,
		&payment.DriverID,
		&payment.Amount,
		&payment.Currency,
		&payment.Status,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		return Payment{}, err
	}

	return payment, nil
}
