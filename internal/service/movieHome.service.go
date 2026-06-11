package service

import (
	"context"
	"log"
	"math"

	"github.com/tickitz-backend/internal/dto"
	"github.com/tickitz-backend/internal/errs"
	"github.com/tickitz-backend/internal/model"

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
			ShowDate:   schedule.ShowDate.Format("2006-01-02"),
			Showtime:   schedule.Showtime.Format("15:04"),
			Price:      schedule.Price,
		})
	}

	return responses, nil
}

func (s *MovieHomeService) GetLocations(
	ctx context.Context,
) ([]dto.MovieLocationRow, error) {
	return s.movieScheduleRepository.FindLocation(ctx)
}

func (s *MovieHomeService) GetShowtimes(
	ctx context.Context,
) ([]dto.MovieShowtimeRow, error) {
	return s.movieScheduleRepository.FindTime(ctx)
}

func (s *MovieHomeService) GetAllMovies(
	ctx context.Context,
	req dto.MovieParamsRequest,
) (*dto.GetAllMoviesResponse, error) {

	if req.Page <= 0 {
		req.Page = 1
	}

	if req.Limit <= 0 {
		req.Limit = 12
	}

	movies, totalData, err := s.movieHomeRepository.GetAllMoviesByFilter(ctx, req)
	if err != nil {
		log.Printf("[MovieHomeService][GetAllMovies] repository error: %v", err)
		return nil, errs.ErrGetMovies
	}

	totalPage := int64(math.Ceil(
		float64(totalData) / float64(req.Limit),
	))

	result := &dto.GetAllMoviesResponse{
		Data: movies,
		Pagination: dto.Meta{
			Page:      req.Page,
			Limit:     req.Limit,
			TotalData: totalData,
			TotalPage: totalPage,
		},
	}

	return result, nil
}

func (s *MovieHomeService) GetUpcomingMovies(
	ctx context.Context,
) ([]model.MoviePreviewResponse, error) {

	movies, err := s.movieHomeRepository.GetUpcomingMovies(ctx)
	if err != nil {
		log.Printf("[MovieHomeService][GetUpcomingMovies] repository error: %v", err)

		return nil, errs.ErrGetMovies
	}

	return movies, nil
}
