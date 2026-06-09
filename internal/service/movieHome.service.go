package service

import (
	"context"

	"github.com/tickitz-backend/internal/dto"
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

func (s *MovieHomeService) GetMovieBySlug(
	ctx context.Context,
	slug string,
	selectedDate string,
	location string,
) (dto.MovieHomeDetailResponse, error) {
	movie, rows, err := s.movieHomeRepo.FindMovieBySlug(
		ctx,
		slug,
		selectedDate,
		location,
	)
	if err != nil {
		return dto.MovieHomeDetailResponse{}, err
	}

	locationMap := make(map[string]map[string][]string)

	for _, row := range rows {
		if _, exists := locationMap[row.Location]; !exists {
			locationMap[row.Location] = make(map[string][]string)
		}

		showtime := row.Showtime.Format("15:04")

		locationMap[row.Location][row.CinemaName] = append(
			locationMap[row.Location][row.CinemaName],
			showtime,
		)
	}

	schedules := make([]dto.LocationScheduleResponse, 0, len(locationMap))

	for locationName, cinemasMap := range locationMap {
		locationResponse := dto.LocationScheduleResponse{
			Location: locationName,
			Cinemas:  make([]dto.CinemaScheduleResponse, 0, len(cinemasMap)),
		}

		for cinemaName, showtimes := range cinemasMap {
			cinemaResponse := dto.CinemaScheduleResponse{
				CinemaName: cinemaName,
				Showtimes:  showtimes,
			}

			locationResponse.Cinemas = append(
				locationResponse.Cinemas,
				cinemaResponse,
			)
		}

		schedules = append(schedules, locationResponse)
	}

	return dto.MovieHomeDetailResponse{
		ID:               movie.ID,
		Title:            movie.Name,
		ReleaseDate:      movie.ReleaseDate.Format("2006-01-02"),
		DurationInMinute: movie.DurationInMinute,
		DirectorName:     movie.DirectorName,
		Synopsis:         movie.Synopsis,
		ImagePoster:      movie.Image,
		GenresCategories: movie.Categories,
		Casts:            movie.Casts,
		Schedules:        schedules,
	}, nil
}
