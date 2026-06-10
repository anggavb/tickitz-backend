package service

import (
	"context"

	"github.com/tickitz-backend/internal/model"
	"github.com/tickitz-backend/internal/repository"
)

type SeatService struct {
	seatRepo *repository.SeatRepository
}

func NewSeatService(seatRepo *repository.SeatRepository) *SeatService {
	return &SeatService{seatRepo: seatRepo}
}

func (s *SeatService) GetSeatMap(ctx context.Context, movieCinemaID int64) (model.SeatMap, error) {
	return s.seatRepo.FindSeatMapByMovieCinemaID(ctx, movieCinemaID)
}
