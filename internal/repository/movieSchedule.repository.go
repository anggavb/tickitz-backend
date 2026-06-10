package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tickitz-backend/internal/dto"
)

type MovieScheduleRepository struct {
	db *pgxpool.Pool
}

func NewMovieScheduleRepository(db *pgxpool.Pool) *MovieScheduleRepository {
	return &MovieScheduleRepository{db: db}
}

func (r *MovieScheduleRepository) FindByMovieSlug(ctx context.Context, slug string) ([]dto.MovieScheduleRow, error) {
	query := `
		SELECT
			l.name AS location,
			c.name AS cinema_name,
			mc.start_date,
			mc.end_date,
			s.showtime,
			mc.price
		FROM movie_cinemas mc
		JOIN movies m ON m.id = mc.movie_id
		JOIN cinemas c ON c.id = mc.cinema_id
		JOIN locations l ON l.id = c.location_id
		JOIN showtimes s ON s.id = mc.showtime_id
		WHERE m.slug = $1
		ORDER BY
			l.name,
			c.name,
			mc.start_date,
			s.showtime;
	`

	rows, err := r.db.Query(ctx, query, slug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	schedules := make([]dto.MovieScheduleRow, 0)

	for rows.Next() {
		var schedule dto.MovieScheduleRow

		err := rows.Scan(
			&schedule.Location,
			&schedule.CinemaName,
			&schedule.StartDate,
			&schedule.EndDate,
			&schedule.Showtime,
			&schedule.Price,
		)
		if err != nil {
			return nil, err
		}

		schedules = append(schedules, schedule)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return schedules, nil
}
