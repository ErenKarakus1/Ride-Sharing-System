package ride

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

func (r *PostgresRepository) Create(ctx context.Context, ride Ride) (Ride, error) {
	const query = `
		INSERT INTO rides (
			id, rider_id, pickup_latitude, pickup_longitude, pickup_address,
			dropoff_latitude, dropoff_longitude, dropoff_address, status
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, rider_id, driver_id, pickup_latitude, pickup_longitude, pickup_address,
			dropoff_latitude, dropoff_longitude, dropoff_address, status, created_at, updated_at
	`

	if ride.ID == "" {
		ride.ID = uuid.NewString()
	}

	return scanRide(r.db.QueryRow(
		ctx,
		query,
		ride.ID,
		strings.TrimSpace(ride.RiderID),
		ride.Pickup.Latitude,
		ride.Pickup.Longitude,
		strings.TrimSpace(ride.Pickup.Address),
		ride.Dropoff.Latitude,
		ride.Dropoff.Longitude,
		strings.TrimSpace(ride.Dropoff.Address),
		StatusRequested,
	))
}

func (r *PostgresRepository) Get(ctx context.Context, id string) (Ride, error) {
	const query = `
		SELECT id, rider_id, driver_id, pickup_latitude, pickup_longitude, pickup_address,
			dropoff_latitude, dropoff_longitude, dropoff_address, status, created_at, updated_at
		FROM rides
		WHERE id = $1
	`

	ride, err := scanRide(r.db.QueryRow(ctx, query, strings.TrimSpace(id)))
	if err == pgx.ErrNoRows {
		return Ride{}, ErrRideNotFound
	}
	if err != nil {
		return Ride{}, err
	}

	return ride, nil
}

func (r *PostgresRepository) ListByRider(ctx context.Context, riderID string) ([]Ride, error) {
	const query = `
		SELECT id, rider_id, driver_id, pickup_latitude, pickup_longitude, pickup_address,
			dropoff_latitude, dropoff_longitude, dropoff_address, status, created_at, updated_at
		FROM rides
		WHERE rider_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, strings.TrimSpace(riderID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rides := make([]Ride, 0)
	for rows.Next() {
		ride, err := scanRide(rows)
		if err != nil {
			return nil, err
		}

		rides = append(rides, ride)
	}

	return rides, rows.Err()
}

func (r *PostgresRepository) UpdateStatus(ctx context.Context, id string, status Status, driverID *string) (Ride, error) {
	const query = `
		UPDATE rides
		SET status = $2, driver_id = COALESCE($3, driver_id), updated_at = now()
		WHERE id = $1
		RETURNING id, rider_id, driver_id, pickup_latitude, pickup_longitude, pickup_address,
			dropoff_latitude, dropoff_longitude, dropoff_address, status, created_at, updated_at
	`

	ride, err := scanRide(r.db.QueryRow(ctx, query, strings.TrimSpace(id), status, driverID))
	if err == pgx.ErrNoRows {
		return Ride{}, ErrRideNotFound
	}
	if err != nil {
		return Ride{}, err
	}

	return ride, nil
}

type rideRow interface {
	Scan(dest ...any) error
}

func scanRide(row rideRow) (Ride, error) {
	var ride Ride
	err := row.Scan(
		&ride.ID,
		&ride.RiderID,
		&ride.DriverID,
		&ride.Pickup.Latitude,
		&ride.Pickup.Longitude,
		&ride.Pickup.Address,
		&ride.Dropoff.Latitude,
		&ride.Dropoff.Longitude,
		&ride.Dropoff.Address,
		&ride.Status,
		&ride.CreatedAt,
		&ride.UpdatedAt,
	)
	if err != nil {
		return Ride{}, err
	}

	return ride, nil
}
