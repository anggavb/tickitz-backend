package repository

import (
	"context"

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

func (r *AuthRepository) Create(ctx context.Context, email string, password string) (int, error) {
	sql := `
	INSERT INTO users (
		email,
		password
	)
	VALUES ($1, $2)
	RETURNING id
	`
	var userID int
	err := r.db.QueryRow(
		ctx,
		sql,
		email,
		password,
	).Scan(&userID)

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
