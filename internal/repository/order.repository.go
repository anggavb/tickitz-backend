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
) ([]dto.OrderHistory, int64, error) {

	offset := (page - 1) * limit

	var totalData int64

	countQuery := `
		SELECT COUNT(*)
		FROM orders
		WHERE user_id = $1
	`

	err := r.db.QueryRow(
		ctx,
		countQuery,
		userID,
	).Scan(&totalData)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT
			o.id,
			o.created_at,

			m.name,
			'' AS movie_category,

			cnm.name,
			cnm.logo,

			mcnm.show_date,
			TO_CHAR(sts.showtime, 'HH24:MI'),

			COALESCE(
				STRING_AGG(
					DISTINCT (s.row || s.number::text),
					', '
				),
				''
			) AS seats,

			COUNT(DISTINCT od.id)::int AS seat_count,

			o.payment_reference,
			o.total_price,

			o.status AS payment_status,

			'' AS ticket_status,

			o.expired_at

		FROM orders o

		JOIN movie_cinemas mcnm
			ON mcnm.id = o.movie_cinema_id

		JOIN movies m
			ON m.id = mcnm.movie_id

		JOIN cinemas cnm
			ON cnm.id = mcnm.cinema_id

		JOIN showtimes sts
			ON sts.id = mcnm.showtime_id

		LEFT JOIN order_details od
			ON od.order_id = o.id

		LEFT JOIN seats s
			ON s.id = od.seat_id

		WHERE o.user_id = $1

		GROUP BY
			o.id,
			o.created_at,
			m.name,
			cnm.name,
			cnm.logo,
			mcnm.show_date,
			sts.showtime,
			o.payment_reference,
			o.total_price,
			o.status,
			o.expired_at

		ORDER BY o.created_at DESC
		LIMIT $2
		OFFSET $3
	`

	rows, err := r.db.Query(
		ctx,
		query,
		userID,
		limit,
		offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	histories := make([]dto.OrderHistory, 0)

	for rows.Next() {
		var history dto.OrderHistory

		err := rows.Scan(
			&history.ID,
			&history.OrderDate,

			&history.MovieName,
			&history.MovieCategory,

			&history.CinemaName,
			&history.CinemaLogo,

			&history.ShowDate,
			&history.ShowTime,

			&history.Seats,
			&history.SeatCount,

			&history.PaymentReference,
			&history.TotalPayment,

			&history.PaymentStatus,
			&history.TicketStatus,

			&history.ExpiredAt,
		)
		if err != nil {
			return nil, 0, err
		}

		histories = append(histories, history)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return histories, totalData, nil
}
