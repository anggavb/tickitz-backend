package service

import (
	"context"
	"encoding/json"
	"math"
	"strconv"
	"strings"

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

func (s *OrderService) GetOrderDetail(ctx context.Context, userID int64, orderID string) (dto.OrderDetailResponse, error) {
	return s.orderRepo.GetOrderDetail(ctx, userID, orderID)
}

func (s *OrderService) UpdateOrderSeats(ctx context.Context, userID int64, orderID string, seats []string) (dto.OrderDetailResponse, error) {
	normalizedSeats := normalizeSeatCodes(seats)
	if len(normalizedSeats) == 0 {
		return dto.OrderDetailResponse{}, repository.ErrInvalidSeats
	}

	return s.orderRepo.UpdateOrderSeats(ctx, userID, orderID, normalizedSeats)
}

func (s *OrderService) GetPaymentMethods(ctx context.Context) ([]dto.OrderPaymentMethodResponse, error) {
	return s.orderRepo.GetPaymentMethods(ctx)
}

func (s *OrderService) SubmitPayment(ctx context.Context, userID int64, orderID string, req dto.UpdateOrderPaymentRequest) (dto.OrderDetailResponse, error) {
	paymentMethodID, err := strconv.ParseInt(req.PaymentMethod, 10, 64)
	if err != nil || paymentMethodID <= 0 {
		return dto.OrderDetailResponse{}, repository.ErrPaymentMethodNotFound
	}

	return s.orderRepo.MarkOrderPaid(ctx, userID, orderID, paymentMethodID)
}

func (s *OrderService) BuildTicketQRPayload(ctx context.Context, userID int64, orderID string) ([]byte, error) {
	order, err := s.orderRepo.GetOrderDetail(ctx, userID, orderID)
	if err != nil {
		return nil, err
	}
	if order.Status != "paid" {
		return nil, repository.ErrOrderNotPayable
	}

	payload := dto.TicketQRPayload{
		OrderID:      order.ID,
		MovieTitle:   order.Movie.Title,
		CinemaName:   order.CinemaName,
		ShowDate:     order.ShowDate,
		ShowTime:     order.ShowTime,
		Seats:        order.Seats,
		SeatCount:    order.SeatCount,
		TotalPayment: order.TotalPayment,
		Status:       order.Status,
	}

	return json.Marshal(payload)
}

func normalizeSeatCodes(seats []string) []string {
	seen := make(map[string]struct{}, len(seats))
	normalized := make([]string, 0, len(seats))

	for _, seat := range seats {
		code := strings.ToUpper(strings.TrimSpace(seat))
		if code == "" {
			continue
		}
		if _, exists := seen[code]; exists {
			continue
		}
		seen[code] = struct{}{}
		normalized = append(normalized, code)
	}

	return normalized
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
