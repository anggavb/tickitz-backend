package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (r *AuthRepository) Create(ctx context.Context, email string, password string, token string) (int64, error) {
	sql := `
		INSERT INTO users (email, password, activation_token, verified_at, token_expiry_at)
			VALUES ($1, $2, $3, NULL, NOW() + INTERVAL '60 minutes')
			RETURNING id
		`
	var userID int64
	err := r.db.QueryRow(ctx, sql, email, password, token).Scan(&userID)

	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (r *AuthRepository) FindByEmail(ctx context.Context, email string) (bool, error) {

	sql := `
	SELECT EXISTS(
		SELECT 1 FROM users WHERE email = $1
	)
	`

	var exists bool

	if err := r.db.QueryRow(
		ctx,
		sql,
		email,
	).Scan(&exists); err != nil {

		return false, err
	}

	return exists, nil
}

func (r *AuthRepository) FindByEmailAndActivate(ctx context.Context, email string) (bool, error) {
	sql := `SELECT EXISTS (
			SELECT 1
			FROM users
			WHERE email = $1
			  AND verified_at IS NULL
		)`

	var exists bool

	if err := r.db.QueryRow(ctx, sql, email).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

func (r *AuthRepository) GetExpiryToken(ctx context.Context, email string) (time.Time, error) {
	sql := `SELECT token_expiry_at FROM users WHERE email = $1;`

	var expiredAt time.Time
	if err := r.db.QueryRow(ctx, sql, email).Scan(&expiredAt); err != nil {
		return time.Time{}, err
	}

	return expiredAt, nil
}

func (r *AuthRepository) GetUserToken(ctx context.Context, email string) (string, error) {
	sql := `SELECT activation_token from users WHERE email = $1`
	var token string
	if err := r.db.QueryRow(ctx, sql, email).Scan(&token); err != nil {
		return "", err
	}
	return token, nil
}

func (r *AuthRepository) Activate(ctx context.Context, email string) error {
	sql := `UPDATE users SET verified_at = NOW() WHERE email = $1`

	_, err := r.db.Exec(ctx, sql, email)
	if err != nil {
		return err
	}
	return nil
}

func (r *AuthRepository) UpdateOTP(ctx context.Context, email string, token string) error {
	sql := `
		UPDATE users
		SET
			activation_token = $1,
			token_expiry_at = NOW() + INTERVAL '60 minutes'
		WHERE email = $2
	`

	_, err := r.db.Exec(ctx, sql, token, email)
	return err
}
