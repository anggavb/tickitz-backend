package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DashboardRepository struct {
	db *pgxpool.Pool
}

type dashboardSalesRow struct {
	Period time.Time
	Total  int64
}

type dashboardTicketRow struct {
	Period time.Time
	Total  int64
	Count  int64
}

func NewDashboardRepository(db *pgxpool.Pool) *DashboardRepository {
	return &DashboardRepository{db: db}
}

var ErrInvalidDashboardPeriod = errors.New("invalid dashboard period")

func (r *DashboardRepository) FindSalesChart(ctx context.Context, movieName, period string) ([]dashboardSalesRow, string, error) {
	trunc, err := normalizeDashboardPeriod(period)
	if err != nil {
		return nil, "", err
	}

	sql := `
SELECT
	date_trunc($1, o.created_at) AS period_date,
	COALESCE(SUM(o.total_price), 0) AS total
FROM orders o
JOIN movie_cinemas mc ON mc.id = o.movie_cinema_id
JOIN movies m ON m.id = mc.movie_id
WHERE o.status = 'paid'
`

	args := []interface{}{trunc}
	if strings.TrimSpace(movieName) != "" {
		sql += ` AND m.name = $2`
		args = append(args, movieName)
	}

	sql += `
GROUP BY period_date
ORDER BY period_date ASC
`

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, trunc, err
	}
	defer rows.Close()

	points := make([]dashboardSalesRow, 0)
	for rows.Next() {
		var row dashboardSalesRow
		if err := rows.Scan(&row.Period, &row.Total); err != nil {
			return nil, trunc, err
		}
		points = append(points, row)
	}

	return points, trunc, rows.Err()
}

func (r *DashboardRepository) FindTicketSales(ctx context.Context, category, location, period string) ([]dashboardTicketRow, string, error) {
	trunc, err := normalizeDashboardPeriod(period)
	if err != nil {
		return nil, "", err
	}

	sql := `
SELECT
	date_trunc($1, o.created_at) AS period_date,
	COUNT(*) AS ticket_count,
	COALESCE(SUM(od.price), 0) AS total
FROM order_details od
JOIN orders o ON o.id = od.order_id
JOIN movie_cinemas mc ON mc.id = o.movie_cinema_id
JOIN movies m ON m.id = mc.movie_id
JOIN cinemas c ON c.id = mc.cinema_id
JOIN locations l ON l.id = c.location_id
WHERE o.status = 'paid'
`

	args := []interface{}{trunc}
	argIndex := 2

	if strings.TrimSpace(category) != "" {
		sql += ` AND EXISTS (
			SELECT 1 FROM movie_categories mcg
			JOIN categories cat ON cat.id = mcg.category_id
			WHERE mcg.movie_id = m.id AND cat.name = $` + strconv.Itoa(argIndex) + `
		)`
		args = append(args, category)
		argIndex++
	}

	if strings.TrimSpace(location) != "" {
		sql += ` AND l.name = $` + strconv.Itoa(argIndex)
		args = append(args, location)
		argIndex++
	}

	sql += `
GROUP BY period_date
ORDER BY period_date ASC
`

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, trunc, err
	}
	defer rows.Close()

	points := make([]dashboardTicketRow, 0)
	for rows.Next() {
		var row dashboardTicketRow
		if err := rows.Scan(&row.Period, &row.Count, &row.Total); err != nil {
			return nil, trunc, err
		}
		points = append(points, row)
	}

	return points, trunc, rows.Err()
}

func normalizeDashboardPeriod(period string) (string, error) {
	period = strings.ToLower(strings.TrimSpace(period))
	if period == "" {
		return "week", nil
	}

	switch period {
	case "daily", "day", "dialy":
		return "day", nil
	case "weekly", "week":
		return "week", nil
	case "monthly", "month":
		return "month", nil
	default:
		return "", fmt.Errorf("%w: %s", ErrInvalidDashboardPeriod, period)
	}
}
