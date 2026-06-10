package service

import (
	"context"
	"strings"
	"time"
	"unicode"

	"github.com/tickitz-backend/internal/dto"
	"github.com/tickitz-backend/internal/model"
	"github.com/tickitz-backend/internal/repository"
)

type MovieService struct {
	movieRepo *repository.MovieRepository
}

type MoviePagination struct {
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
	TotalData int64 `json:"total_data"`
	TotalPage int64 `json:"total_page"`
}

func NewMovieService(movieRepo *repository.MovieRepository) *MovieService {
	return &MovieService{
		movieRepo: movieRepo,
	}
}

func (s *MovieService) List(ctx context.Context, page int, limit int) ([]model.Movie, MoviePagination, error) {
	totalData, err := s.movieRepo.CountAll(ctx)
	if err != nil {
		return nil, MoviePagination{}, err
	}

	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}

	totalPage := int64(0)
	if totalData > 0 {
		totalPage = (totalData + int64(limit) - 1) / int64(limit)
	}

	offset := (page - 1) * limit
	movies, err := s.movieRepo.FindAllPaginated(ctx, limit, offset)
	if err != nil {
		return nil, MoviePagination{}, err
	}

	return movies, MoviePagination{
		Page:      page,
		Limit:     limit,
		TotalData: totalData,
		TotalPage: totalPage,
	}, nil
}

func (s *MovieService) GetByID(ctx context.Context, id int64) (model.Movie, error) {
	return s.movieRepo.FindByID(ctx, id)
}

func (s *MovieService) Create(ctx context.Context, req dto.MovieRequest) (model.Movie, error) {
	releaseDate, err := time.Parse("2006-01-02", req.ReleaseDate)
	if err != nil {
		return model.Movie{}, err
	}

	movie := model.Movie{
		Name:             req.Name,
		Slug:             generateSlug(req.Name),
		ReleaseDate:      releaseDate,
		DurationInMinute: req.DurationInMinute,
		DirectorName:     req.DirectorName,
		Synopsis:         req.Synopsis,
		Image:            req.Image,
		Categories:       req.Categories,
		Casts:            req.Casts,
	}

	movieID, err := s.movieRepo.Create(ctx, movie, req.Categories, req.Casts)
	if err != nil {
		return model.Movie{}, err
	}

	return s.movieRepo.FindByID(ctx, movieID)
}

func (s *MovieService) Update(ctx context.Context, movieID int64, req dto.MovieRequest) (model.Movie, error) {
	releaseDate, err := time.Parse("2006-01-02", req.ReleaseDate)
	if err != nil {
		return model.Movie{}, err
	}

	movie := model.Movie{
		ID:               movieID,
		Name:             req.Name,
		Slug:             generateSlug(req.Name),
		ReleaseDate:      releaseDate,
		DurationInMinute: req.DurationInMinute,
		DirectorName:     req.DirectorName,
		Synopsis:         req.Synopsis,
		Image:            req.Image,
		Categories:       req.Categories,
		Casts:            req.Casts,
	}
	if err := s.movieRepo.Update(ctx, movie, req.Categories, req.Casts); err != nil {
		return model.Movie{}, err
	}
	return s.movieRepo.FindByID(ctx, movieID)
}

func (s *MovieService) Delete(ctx context.Context, movieID int64) error {
	return s.movieRepo.Delete(ctx, movieID)
}

func (s *MovieService) ListCategories(ctx context.Context) ([]string, error) {
	return s.movieRepo.FindAllCategories(ctx)
}

func (s *MovieService) ListCasts(ctx context.Context) ([]string, error) {
	return s.movieRepo.FindAllCasts(ctx)
}

func generateSlug(name string) string {
	name = strings.TrimSpace(strings.ToLower(name))
	var builder strings.Builder
	prevHyphen := false

	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			builder.WriteRune(r)
			prevHyphen = false
			continue
		}
		if !prevHyphen {
			builder.WriteRune('-')
			prevHyphen = true
		}
	}

	slug := strings.Trim(builder.String(), "-")
	return slug
}
