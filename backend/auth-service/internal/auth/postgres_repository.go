package auth

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

func (r *PostgresRepository) Create(ctx context.Context, account Account) (Account, error) {
	const query = `
		INSERT INTO auth_accounts (id, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, email, password_hash, created_at, updated_at
	`

	account.ID = uuid.NewString()
	account.Email = strings.ToLower(strings.TrimSpace(account.Email))

	created, err := scanAccount(r.db.QueryRow(ctx, query, account.ID, account.Email, account.PasswordHash))
	if isUniqueViolation(err) {
		return Account{}, ErrEmailAlreadyExists
	}
	if err != nil {
		return Account{}, err
	}

	return created, nil
}

func (r *PostgresRepository) GetByEmail(ctx context.Context, email string) (Account, error) {
	const query = `
		SELECT id, email, password_hash, created_at, updated_at
		FROM auth_accounts
		WHERE email = $1
	`

	account, err := scanAccount(r.db.QueryRow(ctx, query, strings.ToLower(strings.TrimSpace(email))))
	if err == pgx.ErrNoRows {
		return Account{}, ErrAccountNotFound
	}
	if err != nil {
		return Account{}, err
	}

	return account, nil
}

type accountRow interface {
	Scan(dest ...any) error
}

func scanAccount(row accountRow) (Account, error) {
	var account Account
	err := row.Scan(
		&account.ID,
		&account.Email,
		&account.PasswordHash,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err != nil {
		return Account{}, err
	}

	return account, nil
}

func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}

	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505" ||
		strings.Contains(err.Error(), "SQLSTATE 23505")
}
