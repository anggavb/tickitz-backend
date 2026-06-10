package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tickitz-backend/internal/model"
)

type SeatRepository struct {
	db *pgxpool.Pool
}

func NewSeatRepository(db *pgxpool.Pool) *SeatRepository {
	return &SeatRepository{db: db}
}

func (r *SeatRepository) FindSeatMapByMovieCinemaID(ctx context.Context, movieCinemaID int64) (model.SeatMap, error) {
	seatMap, err := r.findMovieCinemaInfo(ctx, movieCinemaID)
	if err != nil {
		return model.SeatMap{}, err
	}

	seats, err := r.findSeatItemsByMovieCinemaID(ctx, movieCinemaID, seatMap.Cinema.ID)
	if err != nil {
		return model.SeatMap{}, err
	}

	seatMap.Layout = model.SeatLayout{
		Rows:    []string{"A", "B", "C", "D", "E", "F", "G"},
		Columns: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14},
	}
	seatMap.Seats = seats

	return seatMap, nil
}

func (r *SeatRepository) findMovieCinemaInfo(ctx context.Context, movieCinemaID int64) (model.SeatMap, error) {
	query := `
SELECT
	mc.id,
	mc.show_date,
	to_char(st.showtime, 'HH24:MI') AS showtime,
	mc.showtime_id,
	mc.price,
	m.id,
	m.name,
	m.image,
	COALESCE(array_agg(DISTINCT cat.name) FILTER (WHERE cat.name IS NOT NULL), '{}') AS genres,
	c.id,
	c.name,
	COALESCE(c.logo, '') AS logo,
	l.name
FROM movie_cinemas mc
JOIN showtimes st ON st.id = mc.showtime_id
JOIN movies m ON m.id = mc.movie_id
JOIN cinemas c ON c.id = mc.cinema_id
JOIN locations l ON l.id = c.location_id
LEFT JOIN movie_categories mcat ON mcat.movie_id = m.id
LEFT JOIN categories cat ON cat.id = mcat.category_id
WHERE mc.id = $1
GROUP BY mc.id, st.showtime, m.id, c.id, l.id
`
	var seatMap model.SeatMap
	var showDate time.Time
	var genres []string

	err := r.db.QueryRow(ctx, query, movieCinemaID).Scan(
		&seatMap.MovieCinema.ID,
		&showDate,
		&seatMap.MovieCinema.Showtime,
		&seatMap.MovieCinema.ShowtimeID,
		&seatMap.MovieCinema.Price,
		&seatMap.Movie.ID,
		&seatMap.Movie.Title,
		&seatMap.Movie.Poster,
		&genres,
		&seatMap.Cinema.ID,
		&seatMap.Cinema.Name,
		&seatMap.Cinema.Logo,
		&seatMap.Cinema.Location,
	)
	if err != nil {
		return model.SeatMap{}, err
	}

	seatMap.MovieCinema.ShowDate = showDate
	seatMap.Movie.Genres = genres

	return seatMap, nil
}

func (r *SeatRepository) findSeatItemsByMovieCinemaID(ctx context.Context, movieCinemaID int64, cinemaID int64) ([]model.SeatItem, error) {
	query := `
WITH booked_seats AS (
	SELECT od.seat_id
	FROM order_details od
	JOIN orders o ON o.id = od.order_id
	WHERE od.movie_cinema_id = $1
		AND (
			o.status IN ('waiting', 'paid')
			OR (o.status = 'pending' AND o.expired_at > now())
		)
)
SELECT
	s.id,
	trim(s."row") || s."number"::text AS code,
	trim(s."row") AS row_label,
	s."number",
	s."type"::text AS seat_type,
	mc.price,
	CASE
		WHEN bs.seat_id IS NOT NULL THEN 'sold'
		ELSE 'available'
	END AS status
FROM seats s
CROSS JOIN movie_cinemas mc
LEFT JOIN booked_seats bs ON bs.seat_id = s.id
WHERE mc.id = $1
	AND s.cinema_id = $2
ORDER BY s."row", s."number"
`
	rows, err := r.db.Query(ctx, query, movieCinemaID, cinemaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	seats := make([]model.SeatItem, 0, 98)
	for rows.Next() {
		var seat model.SeatItem
		if err := rows.Scan(
			&seat.ID,
			&seat.Code,
			&seat.Row,
			&seat.Number,
			&seat.Type,
			&seat.Price,
			&seat.Status,
		); err != nil {
			return nil, err
		}
		seats = append(seats, seat)
	}

	return seats, rows.Err()
}
