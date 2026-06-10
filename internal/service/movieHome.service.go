package service

import (
	"context"

	"github.com/tickitz-backend/internal/dto"
	"github.com/tickitz-backend/internal/repository"
)

type MovieHomeService struct {
	movieHomeRepository     *repository.MovieHomeRepository
	movieScheduleRepository *repository.MovieScheduleRepository
}

func NewMovieHomeService(
	movieHomeRepository *repository.MovieHomeRepository,
	movieScheduleRepository *repository.MovieScheduleRepository,
) *MovieHomeService {
	return &MovieHomeService{
		movieHomeRepository:     movieHomeRepository,
		movieScheduleRepository: movieScheduleRepository,
	}
}

func (s *MovieHomeService) GetMovieBySlug(ctx context.Context, slug string) (dto.MovieDetailResponse, error) {
	movie, err := s.movieHomeRepository.FindBySlug(ctx, slug)
	if err != nil {
		return dto.MovieDetailResponse{}, err
	}

	return dto.MovieDetailResponse{
		ID:               movie.ID,
		Slug:             movie.Slug,
		Title:            movie.Name,
		ReleaseDate:      movie.ReleaseDate.Format("2006-01-02"),
		DurationInMinute: movie.DurationInMinute,
		DirectorName:     movie.DirectorName,
		Synopsis:         movie.Synopsis,
		ImagePoster:      movie.Image,
		GenresCategories: movie.Categories,
		Casts:            movie.Casts,
	}, nil
}

func (s *MovieHomeService) GetMovieSchedulesBySlug(ctx context.Context, slug string) ([]dto.MovieScheduleResponse, error) {
	_, err := s.movieHomeRepository.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	schedules, err := s.movieScheduleRepository.FindByMovieSlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.MovieScheduleResponse, 0, len(schedules))

	for _, schedule := range schedules {
		responses = append(responses, dto.MovieScheduleResponse{
			Location:   schedule.Location,
			CinemaName: schedule.CinemaName,
			StartDate:  schedule.StartDate.Format("2006-01-02"),
			EndDate:    schedule.EndDate.Format("2006-01-02"),
			Showtime:   schedule.Showtime.Format("15:04"),
			Price:      schedule.Price,
		})
	}

	return responses, nil
}
