package service

import (
	"context"
	"math"

	"github.com/tickitz-backend/internal/dto"
	"github.com/tickitz-backend/internal/repository"
)

type OrderService struct {
	orderRepo *repository.OrderRepository
}

func NewOrderService(orderRepo *repository.OrderRepository) *OrderService {
	return &OrderService{
		orderRepo: orderRepo,
	}
}

func (s *OrderService) CreatePendingOrder(ctx context.Context, userID int64, movieCinemaID int64) (dto.CreatePendingOrderResponse, error) {
	return s.orderRepo.CreatePendingOrder(ctx, userID, movieCinemaID)
}

func (s *OrderService) GetOrderHistory(
	ctx context.Context,
	userID int64,
	page int,
	limit int,
) (*dto.OrderHistoryResponse, error) {

	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 10
	}

	data, totalData, err := s.orderRepo.GetOrderHistory(
		ctx,
		userID,
		page,
		limit,
	)
	if err != nil {
		return nil, err
	}

	totalPage := int(math.Ceil(float64(totalData) / float64(limit)))

	return &dto.OrderHistoryResponse{
		Data: data,
		Meta: dto.Meta{
			Page:      page,
			Limit:     limit,
			TotalData: int64(totalData),
			TotalPage: int64(totalPage),
		},
	}, nil
}
