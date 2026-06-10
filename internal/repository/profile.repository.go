package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tickitz-backend/internal/dto"
	"github.com/tickitz-backend/internal/errs"
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
	sql := `SELECT COALESCE(u.firstname, '') AS first_name, COALESCE(u.lastname, '') AS last_name, COALESCE(u.phone_number, '') AS phone, COALESCE(u.photo, ''), l.name AS loyalty_tier, u.point AS point FROM users u JOIN loyalty_tiers l ON l.id = u.loyalty_tier_id WHERE u.id = $1`
	var profile model.UserProfile
	if err := r.db.QueryRow(ctx, sql, id).Scan(&profile.FirstName, &profile.LastName, &profile.Phone, &profile.Photo, &profile.LoyaltyTier, &profile.Point); err != nil {
		return model.UserProfile{}, err
	}
	return profile, nil
}

func (r *ProfileRepository) UpdateProfile(ctx context.Context, req dto.UpdateProfileRequest, userID int) error {
	var sb strings.Builder
	idx := 1
	args := make([]any, 0)
	sets := make([]string, 0)

	if req.FirstName != nil {
		sets = append(sets, fmt.Sprintf(`first_name = $%d`, idx))
		args = append(args, *req.FirstName)
		idx++
	}
	if req.LastName != nil {
		sets = append(sets, fmt.Sprintf(`last_name = $%d`, idx))
		args = append(args, *req.LastName)
		idx++
	}
	if req.Phone != nil {
		sets = append(sets, fmt.Sprintf(`phone = $%d`, idx))
		args = append(args, *req.Phone)
		idx++
	}
	if req.Photo != nil {
		sets = append(sets, fmt.Sprintf(`photo = $%d`, idx))
		args = append(args, *req.Photo)
		idx++
	}

	if len(sets) == 0 {
		return errs.ErrNothingToUpdate
	}
	sb.WriteString(`UPDATE users SET `)
	sb.WriteString(strings.Join(sets, ", "))
	sb.WriteString(fmt.Sprintf(` WHERE id= $%d`, idx))
	args = append(args, userID)

	query := sb.String()

	if _, err := r.db.Exec(ctx, query, args...); err != nil {
		return err
	}
	return nil

}
