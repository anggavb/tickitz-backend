package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tickitz-backend/internal/dto"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{
		db: db,
	}
}

func (r *OrderRepository) CreatePendingOrder(
	ctx context.Context,
	userID int64,
	movieCinemaID int64,
) (dto.CreatePendingOrderResponse, error) {
	query := `
		WITH selected_schedule AS (
			SELECT
				mc.id AS movie_cinema_id,
				mc.show_date,
				mc.showtime_id,
				m.id AS movie_id,
				m.name AS movie_title,
				m.image AS movie_poster,
				c.id AS cinema_id,
				c.name AS cinema_name,
				COALESCE(c.logo, '') AS cinema_logo,
				l.name AS location,
				to_char(st.showtime, 'HH24:MI') AS showtime
			FROM movie_cinemas mc
			JOIN movies m ON m.id = mc.movie_id
			JOIN cinemas c ON c.id = mc.cinema_id
			JOIN locations l ON l.id = c.location_id
			JOIN showtimes st ON st.id = mc.showtime_id
			WHERE mc.id = $1
				AND (mc.show_date + st.showtime) > now()::timestamp
		),
		existing_order AS (
			SELECT
				o.id::text,
				o.movie_cinema_id,
				o.status::text,
				o.total_price,
				o.expired_at,
				true AS reused
			FROM orders o
			JOIN selected_schedule ss ON ss.movie_cinema_id = o.movie_cinema_id
			WHERE o.user_id = $2
				AND o.movie_cinema_id = $1
				AND o.status = 'pending'
				AND o.expired_at > now()
			ORDER BY o.created_at DESC
			LIMIT 1
		),
		new_order AS (
			INSERT INTO orders (
				user_id,
				showtime_id,
				movie_cinema_id,
				total_price,
				expired_at
			)
			SELECT
				$2,
				showtime_id,
				movie_cinema_id,
				0,
				now() + interval '1 hour'
			FROM selected_schedule
			WHERE NOT EXISTS (SELECT 1 FROM existing_order)
			RETURNING
				id::text,
				movie_cinema_id,
				status::text,
				total_price,
				expired_at,
				false AS reused
		),
		order_result AS (
			SELECT * FROM existing_order
			UNION ALL
			SELECT * FROM new_order
		)
		SELECT
			ord.id,
			ord.movie_cinema_id,
			ord.status,
			ord.total_price,
			ord.expired_at,
			ord.reused,
			ss.movie_id,
			ss.movie_title,
			ss.movie_poster,
			ss.cinema_id,
			ss.cinema_name,
			ss.cinema_logo,
			ss.location,
			to_char(ss.show_date, 'YYYY-MM-DD') AS show_date,
			ss.showtime,
			ss.showtime_id
		FROM order_result ord
		JOIN selected_schedule ss ON ss.movie_cinema_id = ord.movie_cinema_id
	`

	var order dto.CreatePendingOrderResponse
	err := r.db.QueryRow(ctx, query, movieCinemaID, userID).Scan(
		&order.ID,
		&order.MovieCinemaID,
		&order.Status,
		&order.TotalPrice,
		&order.ExpiredAt,
		&order.Reused,
		&order.Movie.ID,
		&order.Movie.Title,
		&order.Movie.Poster,
		&order.Cinema.ID,
		&order.Cinema.Name,
		&order.Cinema.Logo,
		&order.Cinema.Location,
		&order.Schedule.Date,
		&order.Schedule.Time,
		&order.Schedule.ShowtimeID,
	)
	if err != nil {
		return dto.CreatePendingOrderResponse{}, err
	}

	return order, nil
}

func (r *OrderRepository) GetOrderHistory(
	ctx context.Context,
	userID int64,
	page int,
	limit int,
) ([]dto.OrderHistory, int, error) {

	offset := (page - 1) * limit

	// COUNT QUERY
	countQuery := `
		SELECT COUNT(DISTINCT o.id)
		FROM orders o
		WHERE o.user_id = $1
	`

	var total int
	err := r.db.QueryRow(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// MAIN QUERY
	query := `
	SELECT
		o.id,
		o.created_at AS order_date,

		m.name AS movie_name,
		'' AS movie_category,

		cnm.name AS cinema_name,
		COALESCE(cnm.logo, '') AS cinema_logo,

		mc.show_date,
		sts.showtime::text AS show_time,

		COALESCE(STRING_AGG(s.row || s.number::text, ', '), '') AS seats,
		COUNT(od.id)::int AS seat_count,

		COALESCE(o.payment_reference, '') AS payment_reference,
		o.total_price,
		o.status AS payment_status,

		'' AS ticket_status,

		o.expired_at

	FROM orders o
	JOIN movie_cinemas mc ON mc.id = o.movie_cinema_id
	JOIN movies m ON m.id = mc.movie_id
	JOIN cinemas cnm ON cnm.id = mc.cinema_id
	JOIN showtimes sts ON sts.id = o.showtime_id

	LEFT JOIN order_details od ON od.order_id = o.id
	LEFT JOIN seats s ON s.id = od.seat_id

	WHERE o.user_id = $1

	GROUP BY
		o.id, m.id, cnm.id, mc.id, sts.id

	ORDER BY o.created_at DESC
	LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var results []dto.OrderHistory

	for rows.Next() {
		var rdata dto.OrderHistory

		err := rows.Scan(
			&rdata.ID,
			&rdata.OrderDate,
			&rdata.MovieName,
			&rdata.MovieCategory,
			&rdata.CinemaName,
			&rdata.CinemaLogo,
			&rdata.ShowDate,
			&rdata.ShowTime,
			&rdata.Seats,
			&rdata.SeatCount,
			&rdata.PaymentReference,
			&rdata.TotalPayment,
			&rdata.PaymentStatus,
			&rdata.TicketStatus,
			&rdata.ExpiredAt,
		)

		if err != nil {
			return nil, 0, err
		}

		results = append(results, rdata)
	}

	return results, total, nil
}
