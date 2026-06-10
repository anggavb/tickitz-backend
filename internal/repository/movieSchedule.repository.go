package repository

import (
	"context"
	"fmt"

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
            mc.show_date,
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
            mc.show_date,
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
			&schedule.ShowDate,
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

func (r *MovieScheduleRepository) FindLocation(
	ctx context.Context,
) ([]dto.MovieLocationRow, error) {
	query := `
        SELECT DISTINCT name AS location
        FROM locations
        ORDER BY name;
    `
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("FindLocation query: %w", err)
	}
	defer rows.Close()

	locations := make([]dto.MovieLocationRow, 0)
	for rows.Next() {
		var location dto.MovieLocationRow
		if err := rows.Scan(&location.Location); err != nil {
			return nil, fmt.Errorf("FindLocation scan: %w", err)
		}
		locations = append(locations, location)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("FindLocation rows: %w", err)
	}
	return locations, nil
}

func (r *MovieScheduleRepository) FindTime(
	ctx context.Context,
) ([]dto.MovieShowtimeRow, error) {
	query := `
        SELECT DISTINCT showtime
        FROM showtimes
        ORDER BY showtime;
    `
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("FindTime query: %w", err)
	}
	defer rows.Close()

	showtimes := make([]dto.MovieShowtimeRow, 0)
	for rows.Next() {
		var showtime dto.MovieShowtimeRow
		if err := rows.Scan(&showtime.Showtime); err != nil {
			return nil, fmt.Errorf("FindTime scan: %w", err)
		}
		showtimes = append(showtimes, showtime)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("FindTime rows: %w", err)
	}
	return showtimes, nil
}
