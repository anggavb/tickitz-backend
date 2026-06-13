package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tickitz-backend/internal/dto"
)

type MovieScheduleRepository struct {
	db *pgxpool.Pool
}

func NewMovieScheduleRepository(db *pgxpool.Pool) *MovieScheduleRepository {
	return &MovieScheduleRepository{db: db}
}

func (r *MovieScheduleRepository) FindByMovieSlug(ctx context.Context, slug string, filter dto.MovieScheduleQuery) ([]dto.MovieScheduleRow, error) {
	args := []any{slug}
	conditions := []string{"m.slug = $1", "(mc.show_date + s.showtime) > now()::timestamp"}

	if filter.Date != "" {
		args = append(args, filter.Date)
		conditions = append(conditions, fmt.Sprintf("mc.show_date = $%d::date", len(args)))
	}

	if filter.Time != "" {
		args = append(args, filter.Time)
		conditions = append(conditions, fmt.Sprintf("to_char(s.showtime, 'HH24:MI') = $%d", len(args)))
	}

	if filter.Location != "" {
		args = append(args, filter.Location)
		conditions = append(conditions, fmt.Sprintf("l.name = $%d", len(args)))
	}

	query := fmt.Sprintf(`
        SELECT
            mc.id AS movie_cinema_id,
            c.id AS cinema_id,
            l.name AS location,
            c.name AS cinema_name,
            COALESCE(c.logo, '') AS cinema_logo,
            mc.show_date,
            s.showtime,
            s.id AS showtime_id,
            mc.price
        FROM movie_cinemas mc
        JOIN movies m ON m.id = mc.movie_id
        JOIN cinemas c ON c.id = mc.cinema_id
        JOIN locations l ON l.id = c.location_id
        JOIN showtimes s ON s.id = mc.showtime_id
        WHERE %s
        ORDER BY
            l.name,
            c.name,
            mc.show_date,
            s.showtime;
    `, strings.Join(conditions, " AND "))

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	schedules := make([]dto.MovieScheduleRow, 0)

	for rows.Next() {
		var schedule dto.MovieScheduleRow

		err := rows.Scan(
			&schedule.MovieCinemaID,
			&schedule.CinemaID,
			&schedule.Location,
			&schedule.CinemaName,
			&schedule.CinemaLogo,
			&schedule.ShowDate,
			&schedule.Showtime,
			&schedule.ShowtimeID,
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

func (r *MovieScheduleRepository) FindOptionsByMovieSlug(ctx context.Context, slug string) (dto.MovieScheduleOptionsResponse, error) {
	dates, err := r.findDatesByMovieSlug(ctx, slug)
	if err != nil {
		return dto.MovieScheduleOptionsResponse{}, err
	}

	showtimes, err := r.findShowtimesByMovieSlug(ctx, slug)
	if err != nil {
		return dto.MovieScheduleOptionsResponse{}, err
	}

	locations, err := r.findLocationsByMovieSlug(ctx, slug)
	if err != nil {
		return dto.MovieScheduleOptionsResponse{}, err
	}

	return dto.MovieScheduleOptionsResponse{
		Dates:     dates,
		Showtimes: showtimes,
		Locations: locations,
	}, nil
}

func (r *MovieScheduleRepository) findDatesByMovieSlug(ctx context.Context, slug string) ([]string, error) {
	query := `
        SELECT DISTINCT to_char(mc.show_date, 'YYYY-MM-DD') AS show_date
        FROM movie_cinemas mc
        JOIN movies m ON m.id = mc.movie_id
        JOIN showtimes s ON s.id = mc.showtime_id
        WHERE m.slug = $1
            AND (mc.show_date + s.showtime) > now()::timestamp
        ORDER BY show_date;
    `

	rows, err := r.db.Query(ctx, query, slug)
	if err != nil {
		return nil, fmt.Errorf("findDatesByMovieSlug query: %w", err)
	}
	defer rows.Close()

	dates := make([]string, 0)
	for rows.Next() {
		var date string
		if err := rows.Scan(&date); err != nil {
			return nil, fmt.Errorf("findDatesByMovieSlug scan: %w", err)
		}
		dates = append(dates, date)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("findDatesByMovieSlug rows: %w", err)
	}

	return dates, nil
}

func (r *MovieScheduleRepository) findShowtimesByMovieSlug(ctx context.Context, slug string) ([]dto.MovieShowtimeRow, error) {
	query := `
        SELECT DISTINCT to_char(s.showtime, 'HH24:MI') AS showtime
        FROM movie_cinemas mc
        JOIN movies m ON m.id = mc.movie_id
        JOIN showtimes s ON s.id = mc.showtime_id
        WHERE m.slug = $1
            AND (mc.show_date + s.showtime) > now()::timestamp
        ORDER BY showtime;
    `

	rows, err := r.db.Query(ctx, query, slug)
	if err != nil {
		return nil, fmt.Errorf("findShowtimesByMovieSlug query: %w", err)
	}
	defer rows.Close()

	showtimes := make([]dto.MovieShowtimeRow, 0)
	for rows.Next() {
		var showtime dto.MovieShowtimeRow
		if err := rows.Scan(&showtime.Showtime); err != nil {
			return nil, fmt.Errorf("findShowtimesByMovieSlug scan: %w", err)
		}
		showtimes = append(showtimes, showtime)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("findShowtimesByMovieSlug rows: %w", err)
	}

	return showtimes, nil
}

func (r *MovieScheduleRepository) findLocationsByMovieSlug(ctx context.Context, slug string) ([]dto.MovieLocationRow, error) {
	query := `
        SELECT DISTINCT l.name AS location
        FROM movie_cinemas mc
        JOIN movies m ON m.id = mc.movie_id
        JOIN cinemas c ON c.id = mc.cinema_id
        JOIN locations l ON l.id = c.location_id
        JOIN showtimes s ON s.id = mc.showtime_id
        WHERE m.slug = $1
            AND (mc.show_date + s.showtime) > now()::timestamp
        ORDER BY l.name;
    `

	rows, err := r.db.Query(ctx, query, slug)
	if err != nil {
		return nil, fmt.Errorf("findLocationsByMovieSlug query: %w", err)
	}
	defer rows.Close()

	locations := make([]dto.MovieLocationRow, 0)
	for rows.Next() {
		var location dto.MovieLocationRow
		if err := rows.Scan(&location.Location); err != nil {
			return nil, fmt.Errorf("findLocationsByMovieSlug scan: %w", err)
		}
		locations = append(locations, location)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("findLocationsByMovieSlug rows: %w", err)
	}

	return locations, nil
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
        SELECT DISTINCT to_char(showtime, 'HH24:MI') AS showtime
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

func (r *MovieScheduleRepository) FindAllCinemas(ctx context.Context) ([]dto.CinemaResponse, error) {
	query := `
        SELECT c.id, c.name, COALESCE(l.name, '') AS location
        FROM cinemas c
        LEFT JOIN locations l ON l.id = c.location_id
        ORDER BY c.name;
    `

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("FindAllCinemas query: %w", err)
	}
	defer rows.Close()

	cinemas := make([]dto.CinemaResponse, 0)
	for rows.Next() {
		var cinema dto.CinemaResponse
		if err := rows.Scan(&cinema.ID, &cinema.Name, &cinema.Location); err != nil {
			return nil, fmt.Errorf("FindAllCinemas scan: %w", err)
		}
		cinemas = append(cinemas, cinema)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("FindAllCinemas rows: %w", err)
	}

	return cinemas, nil
}

func (r *MovieScheduleRepository) FindShowtimesByMovieID(ctx context.Context, movieID int64) ([]dto.MovieScheduleResponse, error) {
	query := `
        SELECT
            mc.id AS movie_cinema_id,
            c.id AS cinema_id,
            l.name AS location,
            c.name AS cinema_name,
            COALESCE(c.logo, '') AS cinema_logo,
            mc.show_date,
            s.showtime,
            s.id AS showtime_id,
            mc.price
        FROM movie_cinemas mc
        JOIN cinemas c ON c.id = mc.cinema_id
        LEFT JOIN locations l ON l.id = c.location_id
        JOIN showtimes s ON s.id = mc.showtime_id
        WHERE mc.movie_id = $1
        ORDER BY l.name, c.name, mc.show_date, s.showtime;
    `

	rows, err := r.db.Query(ctx, query, movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	schedules := make([]dto.MovieScheduleResponse, 0)
	for rows.Next() {
		var row dto.MovieScheduleRow
		if err := rows.Scan(
			&row.MovieCinemaID,
			&row.CinemaID,
			&row.Location,
			&row.CinemaName,
			&row.CinemaLogo,
			&row.ShowDate,
			&row.Showtime,
			&row.ShowtimeID,
			&row.Price,
		); err != nil {
			return nil, err
		}
		schedules = append(schedules, dto.MovieScheduleResponse{
			MovieCinemaID: row.MovieCinemaID,
			CinemaID:      row.CinemaID,
			Location:      row.Location,
			CinemaName:    row.CinemaName,
			CinemaLogo:    row.CinemaLogo,
			ShowDate:      row.ShowDate.Format("2006-01-02"),
			Showtime:      row.Showtime.Format("15:04"),
			ShowtimeID:    row.ShowtimeID,
			Price:         row.Price,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return schedules, nil
}

func (r *MovieScheduleRepository) findShowtimeIDsByTimes(ctx context.Context, times []string) (map[string]int64, error) {
	query := `
        SELECT id, to_char(showtime, 'HH24:MI') AS showtime
        FROM showtimes
        WHERE to_char(showtime, 'HH24:MI') = ANY($1::text[])
    `

	rows, err := r.db.Query(ctx, query, times)
	if err != nil {
		return nil, fmt.Errorf("findShowtimeIDsByTimes query: %w", err)
	}
	defer rows.Close()

	showtimeIDs := make(map[string]int64)
	for rows.Next() {
		var id int64
		var showtime string
		if err := rows.Scan(&id, &showtime); err != nil {
			return nil, fmt.Errorf("findShowtimeIDsByTimes scan: %w", err)
		}
		showtimeIDs[showtime] = id
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("findShowtimeIDsByTimes rows: %w", err)
	}

	return showtimeIDs, nil
}

func (r *MovieScheduleRepository) UpsertMovieCinemas(ctx context.Context, movieID int64, cinemaID int64, startDate string, endDate string, times []string, price int) error {
	if len(times) == 0 {
		return fmt.Errorf("at least one showtime is required")
	}

	distinctTimes := make([]string, 0, len(times))
	timeSet := make(map[string]struct{})
	for _, t := range times {
		if t == "" {
			continue
		}
		if _, exists := timeSet[t]; !exists {
			timeSet[t] = struct{}{}
			distinctTimes = append(distinctTimes, t)
		}
	}

	if len(distinctTimes) == 0 {
		return fmt.Errorf("at least one showtime is required")
	}

	showtimeIDs, err := r.findShowtimeIDsByTimes(ctx, distinctTimes)
	if err != nil {
		return err
	}

	if len(showtimeIDs) != len(distinctTimes) {
		missing := make([]string, 0)
		for _, t := range distinctTimes {
			if _, ok := showtimeIDs[t]; !ok {
				missing = append(missing, t)
			}
		}
		return fmt.Errorf("invalid showtimes: %v", missing)
	}

	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	insertSQL := `
        INSERT INTO movie_cinemas (movie_id, cinema_id, show_date, showtime_id, price)
        SELECT $1, $2, generated_date::date, st.id, $5
        FROM generate_series($3::date, $4::date, interval '1 day') AS generated_date
        JOIN showtimes st ON to_char(st.showtime, 'HH24:MI') = ANY($6::text[])
        ON CONFLICT (movie_id, cinema_id, show_date, showtime_id) DO UPDATE
            SET price = EXCLUDED.price;
    `

	_, err = tx.Exec(ctx, insertSQL, movieID, cinemaID, startDate, endDate, price, distinctTimes)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}
