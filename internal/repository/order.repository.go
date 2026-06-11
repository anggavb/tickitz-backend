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
		cnm.logo AS cinema_logo,

		mc.show_date,
		sts.showtime::text AS show_time,

		COALESCE(STRING_AGG(s.row || s.number::text, ', '), '') AS seats,
		COUNT(od.id)::int AS seat_count,

		o.payment_reference,
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
