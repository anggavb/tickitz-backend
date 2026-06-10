package service

import (
	"context"

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

func (s *MovieHomeService) GetAllMovies(ctx context.Context, req dto.MovieParamsRequest) ([]model.MoviePreviewResponse, error) {
	return s.movieHomeRepo.GetAllMoviesByFilter(ctx, req)
}
