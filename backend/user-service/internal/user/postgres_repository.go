package user

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) Create(ctx context.Context, user User) (User, error) {
	const query = `
		INSERT INTO users (id, email, display_name, phone_number, role)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, email, display_name, phone_number, role, created_at, updated_at
	`

	user.ID = uuid.NewString()

	row := r.db.QueryRow(
		ctx,
		query,
		user.ID,
		strings.ToLower(strings.TrimSpace(user.Email)),
		strings.TrimSpace(user.DisplayName),
		strings.TrimSpace(user.PhoneNumber),
		user.Role,
	)

	created, err := scanUser(row)
	if isUniqueViolation(err) {
		return User{}, ErrEmailAlreadyExists
	}
	if err != nil {
		return User{}, err
	}

	return created, nil
}

func (r *PostgresRepository) List(ctx context.Context) ([]User, error) {
	const query = `
		SELECT id, email, display_name, phone_number, role, created_at, updated_at
		FROM users
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, rows.Err()
}

func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}

	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505" ||
		strings.Contains(err.Error(), "SQLSTATE 23505")
}

func (r *PostgresRepository) Get(ctx context.Context, id string) (User, error) {
	const query = `
		SELECT id, email, display_name, phone_number, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user, err := scanUser(r.db.QueryRow(ctx, query, id))
	if err == pgx.ErrNoRows {
		return User{}, ErrUserNotFound
	}
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (r *PostgresRepository) Update(ctx context.Context, id string, request UpdateUserRequest) (User, error) {
	const query = `
		UPDATE users
		SET
			display_name = COALESCE($2, display_name),
			phone_number = COALESCE($3, phone_number),
			updated_at = now()
		WHERE id = $1
		RETURNING id, email, display_name, phone_number, role, created_at, updated_at
	`

	var displayName *string
	if request.DisplayName != nil {
		trimmed := strings.TrimSpace(*request.DisplayName)
		displayName = &trimmed
	}

	var phoneNumber *string
	if request.PhoneNumber != nil {
		trimmed := strings.TrimSpace(*request.PhoneNumber)
		phoneNumber = &trimmed
	}

	user, err := scanUser(r.db.QueryRow(ctx, query, id, displayName, phoneNumber))
	if err == pgx.ErrNoRows {
		return User{}, ErrUserNotFound
	}
	if err != nil {
		return User{}, err
	}

	return user, nil
}

type userRow interface {
	Scan(dest ...any) error
}

func scanUser(row userRow) (User, error) {
	var user User
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.DisplayName,
		&user.PhoneNumber,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
