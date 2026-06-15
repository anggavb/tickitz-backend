package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tickitz-backend/internal/dto"
)

var (
	ErrOrderNotFound         = errors.New("order not found")
	ErrOrderExpired          = errors.New("order expired")
	ErrOrderNotPayable       = errors.New("order is not ready for payment")
	ErrOrderAlreadyPaid      = errors.New("order already paid")
	ErrInvalidSeats          = errors.New("invalid seats")
	ErrSeatUnavailable       = errors.New("seat unavailable")
	ErrPaymentMethodNotFound = errors.New("payment method not found")
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

func (r *OrderRepository) GetOrderDetail(ctx context.Context, userID int64, orderID string) (dto.OrderDetailResponse, error) {
	query := `
SELECT
	o.id::text,
	o.status::text,
	o.movie_cinema_id,
	to_char(mc.show_date, 'YYYY-MM-DD') AS show_date,
	to_char(st.showtime, 'HH24:MI') AS show_time,
	l.name AS location,
	c.id AS cinema_id,
	c.name AS cinema_name,
	COALESCE(c.logo, '') AS cinema_logo,
	mc.price AS ticket_price,
	COALESCE(seat_data.seats, '{}') AS seats,
	COALESCE(array_length(seat_data.seats, 1), 0) AS seat_count,
	o.total_price,
	COALESCE(o.payment_reference, '') AS payment_reference,
	o.expired_at,
	m.id AS movie_id,
	m.name AS movie_title,
	COALESCE(m.image, '') AS movie_poster,
	COALESCE(category_data.genres, '{}') AS genres,
	COALESCE(pm.id, 0) AS payment_method_id,
	COALESCE(pm.name, '') AS payment_method_name,
	COALESCE(pm.logo, '') AS payment_method_logo
FROM orders o
JOIN movie_cinemas mc ON mc.id = o.movie_cinema_id
JOIN movies m ON m.id = mc.movie_id
JOIN cinemas c ON c.id = mc.cinema_id
JOIN locations l ON l.id = c.location_id
JOIN showtimes st ON st.id = o.showtime_id
LEFT JOIN payment_methods pm ON pm.id = o.payment_method_id
LEFT JOIN LATERAL (
	SELECT array_agg(trim(s."row") || s."number"::text ORDER BY s."row", s."number") AS seats
	FROM order_details od
	JOIN seats s ON s.id = od.seat_id
	WHERE od.order_id = o.id
) seat_data ON true
LEFT JOIN LATERAL (
	SELECT array_agg(DISTINCT cat.name ORDER BY cat.name) AS genres
	FROM movie_categories mcat
	JOIN categories cat ON cat.id = mcat.category_id
	WHERE mcat.movie_id = m.id
) category_data ON true
WHERE o.id = $1
	AND o.user_id = $2
`

	var order dto.OrderDetailResponse
	var paymentMethodID int64
	var paymentMethodName string
	var paymentMethodLogo string

	err := r.db.QueryRow(ctx, query, orderID, userID).Scan(
		&order.ID,
		&order.Status,
		&order.MovieCinemaID,
		&order.ShowDate,
		&order.ShowTime,
		&order.Location,
		&order.Cinema.ID,
		&order.Cinema.Name,
		&order.Cinema.Logo,
		&order.TicketPrice,
		&order.Seats,
		&order.SeatCount,
		&order.TotalPayment,
		&order.PaymentReference,
		&order.ExpiredAt,
		&order.Movie.ID,
		&order.Movie.Title,
		&order.Movie.Poster,
		&order.Movie.Genres,
		&paymentMethodID,
		&paymentMethodName,
		&paymentMethodLogo,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.OrderDetailResponse{}, ErrOrderNotFound
		}
		return dto.OrderDetailResponse{}, err
	}

	order.Movie.Background = order.Movie.Poster
	order.CinemaName = order.Cinema.Name
	order.PaymentStatus = order.Status
	if order.Status == "paid" {
		order.TicketStatus = "active"
	} else {
		order.TicketStatus = order.Status
	}
	if paymentMethodID > 0 {
		order.PaymentMethod = &dto.OrderPaymentMethodResponse{
			ID:    paymentMethodID,
			Name:  paymentMethodName,
			Label: paymentMethodName,
			Logo:  paymentMethodLogo,
		}
	}

	return order, nil
}

func (r *OrderRepository) UpdateOrderSeats(ctx context.Context, userID int64, orderID string, seatCodes []string) (dto.OrderDetailResponse, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return dto.OrderDetailResponse{}, err
	}
	defer tx.Rollback(ctx)

	var movieCinemaID int64
	var showtimeID int64
	var cinemaID int64
	var ticketPrice int
	var status string

	orderQuery := `
SELECT o.movie_cinema_id, o.showtime_id, mc.cinema_id, mc.price, o.status::text
FROM orders o
JOIN movie_cinemas mc ON mc.id = o.movie_cinema_id
WHERE o.id = $1
	AND o.user_id = $2
FOR UPDATE
`
	err = tx.QueryRow(ctx, orderQuery, orderID, userID).Scan(
		&movieCinemaID,
		&showtimeID,
		&cinemaID,
		&ticketPrice,
		&status,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.OrderDetailResponse{}, ErrOrderNotFound
		}
		return dto.OrderDetailResponse{}, err
	}

	if status == "paid" {
		return dto.OrderDetailResponse{}, ErrOrderAlreadyPaid
	}
	if status == "cancel" {
		return dto.OrderDetailResponse{}, ErrOrderNotFound
	}

	rows, err := tx.Query(ctx, `
SELECT id, trim("row") || "number"::text AS code
FROM seats
WHERE cinema_id = $1
	AND trim("row") || "number"::text = ANY($2)
`, cinemaID, seatCodes)
	if err != nil {
		return dto.OrderDetailResponse{}, err
	}
	defer rows.Close()

	seatIDsByCode := make(map[string]int64, len(seatCodes))
	for rows.Next() {
		var seatID int64
		var code string
		if err := rows.Scan(&seatID, &code); err != nil {
			return dto.OrderDetailResponse{}, err
		}
		seatIDsByCode[code] = seatID
	}
	if err := rows.Err(); err != nil {
		return dto.OrderDetailResponse{}, err
	}
	if len(seatIDsByCode) != len(seatCodes) {
		return dto.OrderDetailResponse{}, ErrInvalidSeats
	}

	if _, err = tx.Exec(ctx, `DELETE FROM order_details WHERE order_id = $1`, orderID); err != nil {
		return dto.OrderDetailResponse{}, err
	}

	for _, seatCode := range seatCodes {
		seatID := seatIDsByCode[seatCode]
		_, err = tx.Exec(ctx, `
INSERT INTO order_details (order_id, seat_id, showtime_id, movie_cinema_id, price)
VALUES ($1, $2, $3, $4, $5)
`, orderID, seatID, showtimeID, movieCinemaID, ticketPrice)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				return dto.OrderDetailResponse{}, ErrSeatUnavailable
			}
			return dto.OrderDetailResponse{}, err
		}
	}

	totalPrice := ticketPrice * len(seatCodes)
	_, err = tx.Exec(ctx, `
UPDATE orders
SET total_price = $1,
	status = 'waiting'
WHERE id = $2
`, totalPrice, orderID)
	if err != nil {
		return dto.OrderDetailResponse{}, err
	}

	_, err = tx.Exec(ctx, `
WITH order_points AS (
	SELECT (total_price / 1000)::int AS earned_points
	FROM orders
	WHERE id = $1
		AND user_id = $2
),
new_points AS (
	SELECT u.point + op.earned_points AS point
	FROM users u
	CROSS JOIN order_points op
	WHERE u.id = $2
),
new_tier AS (
	SELECT lt.id
	FROM loyalty_tiers lt
	JOIN new_points np ON lt.min_point <= np.point
	ORDER BY lt.min_point DESC
	LIMIT 1
)
UPDATE users
SET point = (SELECT point FROM new_points),
	loyalty_tier_id = (SELECT id FROM new_tier),
	updated_at = now()
WHERE id = $2
`, orderID, userID)
	if err != nil {
		return dto.OrderDetailResponse{}, err
	}

	if err = tx.Commit(ctx); err != nil {
		return dto.OrderDetailResponse{}, err
	}

	return r.GetOrderDetail(ctx, userID, orderID)
}

func (r *OrderRepository) GetPaymentMethods(ctx context.Context) ([]dto.OrderPaymentMethodResponse, error) {
	rows, err := r.db.Query(ctx, `
SELECT id, name, logo
FROM payment_methods
ORDER BY id
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	methods := make([]dto.OrderPaymentMethodResponse, 0)
	for rows.Next() {
		var method dto.OrderPaymentMethodResponse
		if err := rows.Scan(&method.ID, &method.Name, &method.Logo); err != nil {
			return nil, err
		}
		method.Label = method.Name
		methods = append(methods, method)
	}

	return methods, rows.Err()
}

func (r *OrderRepository) MarkOrderPaid(ctx context.Context, userID int64, orderID string, paymentMethodID int64) (dto.OrderDetailResponse, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return dto.OrderDetailResponse{}, err
	}
	defer tx.Rollback(ctx)

	var methodExists bool
	if err = tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM payment_methods WHERE id = $1)`, paymentMethodID).Scan(&methodExists); err != nil {
		return dto.OrderDetailResponse{}, err
	}
	if !methodExists {
		return dto.OrderDetailResponse{}, ErrPaymentMethodNotFound
	}

	var status string
	var seatCount int
	err = tx.QueryRow(ctx, `
SELECT status::text
FROM orders
WHERE id = $1
	AND user_id = $2
FOR UPDATE
`, orderID, userID).Scan(&status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.OrderDetailResponse{}, ErrOrderNotFound
		}
		return dto.OrderDetailResponse{}, err
	}
	if err = tx.QueryRow(ctx, `SELECT COUNT(*)::int FROM order_details WHERE order_id = $1`, orderID).Scan(&seatCount); err != nil {
		return dto.OrderDetailResponse{}, err
	}
	if status == "paid" {
		return dto.OrderDetailResponse{}, ErrOrderAlreadyPaid
	}
	if status != "waiting" || seatCount == 0 {
		return dto.OrderDetailResponse{}, ErrOrderNotPayable
	}

	paymentReference := fmt.Sprintf("TICKITZ-%s", strings.ToUpper(strings.ReplaceAll(orderID, "-", ""))[:12])
	_, err = tx.Exec(ctx, `
UPDATE orders
SET payment_method_id = $1,
	payment_reference = $2,
	status = 'paid'
WHERE id = $3
`, paymentMethodID, paymentReference, orderID)
	if err != nil {
		return dto.OrderDetailResponse{}, err
	}

	if err = tx.Commit(ctx); err != nil {
		return dto.OrderDetailResponse{}, err
	}

	return r.GetOrderDetail(ctx, userID, orderID)
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
