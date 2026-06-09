package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tickitz-backend/internal/model"
)

type ProfileRepository struct {
	db *pgxpool.Pool
}

func NewProfileRepository(db *pgxpool.Pool) *ProfileRepository {
	return &ProfileRepository{
		db: db,
	}
}

func (r *ProfileRepository) GetProfile(ctx context.Context, id int) (model.UserProfile, error) {
	sql := `SELECT COALESCE(u.firstname, '') AS first_name, COALESCE(u.lastname, '') AS last_name, u.phone_number AS phone, u.photo, l.name AS loyalty_tier FROM users u JOIN loyalty_tiers l ON l.id = u.loyalty_tier_id WHERE u.id = $1`
	var profile model.UserProfile
	if err := r.db.QueryRow(ctx, sql, id).Scan(&profile.FirstName, &profile.LastName, &profile.Phone, &profile.Photo, &profile.LoyaltyTier); err != nil {
		return model.UserProfile{}, err
	}
	return profile, nil
}
