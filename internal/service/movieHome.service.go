package service

import (
	"context"
	"errors"

	"github.com/tickitz-backend/internal/dto"
	"github.com/tickitz-backend/internal/model"

	"github.com/tickitz-backend/internal/repository"
)

type MovieHomeService struct {
	movieHomeRepo *repository.MovieHomeRepository
}

func NewMovieHomeService(movieHomeRepo *repository.MovieHomeRepository) *MovieHomeService {
	return &MovieHomeService{
		movieHomeRepo: movieHomeRepo,
	}
}

func (s *MovieHomeService) GetBySlug(ctx context.Context, slug string) (model.MovieDetails, error) {
	return s.movieHomeRepo.FindBySlug(ctx, slug)
}
func (s *MovieHomeService) GetScheduleBySlugAndLocation(ctx context.Context, slug string, location string) ([]dto.LocationScheduleResponse, error) {
	if slug == "" {
		return nil, errors.New("slug path parameter is required")
	}

	schedules, err := s.movieHomeRepo.FindScheduleBySlugAndLocation(ctx, slug, location)
	if err != nil {
		return nil, err
	}

	if schedules == nil {
		return []dto.LocationScheduleResponse{}, nil
	}

	return schedules, nil
}
