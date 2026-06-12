package service

import (
	"context"
	"log"
	"math"
	"time"

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

func (s *MovieHomeService) GetMovieSchedulesBySlug(ctx context.Context, slug string, filter dto.MovieScheduleQuery) ([]dto.MovieScheduleResponse, error) {
	if err := validateScheduleFilter(filter); err != nil {
		return nil, err
	}

	_, err := s.movieHomeRepository.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	schedules, err := s.movieScheduleRepository.FindByMovieSlug(ctx, slug, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.MovieScheduleResponse, 0, len(schedules))

	for _, schedule := range schedules {
		responses = append(responses, dto.MovieScheduleResponse{
			MovieCinemaID: schedule.MovieCinemaID,
			CinemaID:      schedule.CinemaID,
			Location:      schedule.Location,
			CinemaName:    schedule.CinemaName,
			CinemaLogo:    schedule.CinemaLogo,
			ShowDate:      schedule.ShowDate.Format("2006-01-02"),
			Showtime:      schedule.Showtime.Format("15:04"),
			ShowtimeID:    schedule.ShowtimeID,
			Price:         schedule.Price,
		})
	}

	return responses, nil
}

func validateScheduleFilter(filter dto.MovieScheduleQuery) error {
	now := time.Now()
	today := now.Format("2006-01-02")

	if filter.Date != "" {
		parsedDate, err := time.Parse("2006-01-02", filter.Date)
		if err != nil {
			return errs.ErrInvalidScheduleFilter
		}

		currentDate, _ := time.Parse("2006-01-02", today)
		if parsedDate.Before(currentDate) {
			return errs.ErrInvalidScheduleFilter
		}
	}

	if filter.Time == "" {
		return nil
	}

	parsedTime, err := time.Parse("15:04", filter.Time)
	if err != nil {
		return errs.ErrInvalidScheduleFilter
	}

	// A time-only filter is interpreted against today because no future date
	// context exists yet. Future dates may still use any showtime.
	if filter.Date != "" && filter.Date != today {
		return nil
	}

	filterMinutes := parsedTime.Hour()*60 + parsedTime.Minute()
	nowMinutes := now.Hour()*60 + now.Minute()
	if filterMinutes <= nowMinutes {
		return errs.ErrInvalidScheduleFilter
	}

	return nil
}

func (s *MovieHomeService) GetMovieScheduleOptionsBySlug(ctx context.Context, slug string) (dto.MovieScheduleOptionsResponse, error) {
	_, err := s.movieHomeRepository.FindBySlug(ctx, slug)
	if err != nil {
		return dto.MovieScheduleOptionsResponse{}, err
	}

	return s.movieScheduleRepository.FindOptionsByMovieSlug(ctx, slug)
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
