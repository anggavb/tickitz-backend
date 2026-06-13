package service

import (
	"context"
	"time"

	"github.com/tickitz-backend/internal/dto"
	"github.com/tickitz-backend/internal/repository"
)

type DashboardService struct {
	dashboardRepo *repository.DashboardRepository
}

func NewDashboardService(dashboardRepo *repository.DashboardRepository) *DashboardService {
	return &DashboardService{dashboardRepo: dashboardRepo}
}

func (s *DashboardService) GetSalesChart(ctx context.Context, movieName, period string) ([]dto.DashboardSalesPoint, string, error) {
	rows, normalizedPeriod, err := s.dashboardRepo.FindSalesChart(ctx, movieName, period)
	if err != nil {
		return nil, "", err
	}

	points := make([]dto.DashboardSalesPoint, 0, len(rows))
	for _, row := range rows {
		points = append(points, dto.DashboardSalesPoint{
			Period:  formatDashboardPeriodLabel(normalizedPeriod, row.Period),
			Revenue: row.Total,
		})
	}

	return points, normalizedPeriod, nil
}

func (s *DashboardService) GetTicketSales(ctx context.Context, category, location, period string) ([]dto.DashboardTicketSalesPoint, string, error) {
	rows, normalizedPeriod, err := s.dashboardRepo.FindTicketSales(ctx, category, location, period)
	if err != nil {
		return nil, "", err
	}

	points := make([]dto.DashboardTicketSalesPoint, 0, len(rows))
	for _, row := range rows {
		points = append(points, dto.DashboardTicketSalesPoint{
			Period:      formatDashboardPeriodLabel(normalizedPeriod, row.Period),
			TicketCount: row.Count,
			Revenue:     row.Total,
		})
	}

	return points, normalizedPeriod, nil
}

func formatDashboardPeriodLabel(period string, timestamp time.Time) string {
	switch period {
	case "day":
		return timestamp.Format("2006-01-02")
	case "month":
		return timestamp.Format("2006-01")
	default:
		return timestamp.Format("2006-01-02")
	}
}
