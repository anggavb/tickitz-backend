package controller

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/tickitz-backend/internal/dto"
	"github.com/tickitz-backend/internal/service"
)

type MovieHomeController struct {
	movieHomeService *service.MovieHomeService
}

func NewMovieHomeController(movieHomeService *service.MovieHomeService) *MovieHomeController {
	return &MovieHomeController{
		movieHomeService: movieHomeService,
	}
}

// GetBySlug godoc
// @Summary      Get a Movie by its slug
// @Description  Retrieve detailed information about a specific movie using its unique URI slug.
// @Tags         Movies
// @Accept       json
// @Produce      json
// @Param        slug  path      string  true  "Movie Slug"
// @Success      200   {object}  dto.MovieDetailWrappedResponse
// @Failure      400   {object}  dto.ErrorResponse
// @Failure      404   {object}  dto.ErrorResponse
// @Failure      500   {object}  dto.ErrorResponse
// @Router       /movies/{slug} [get]
func (c *MovieHomeController) GetBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")
	if slug == "" {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Status:  "error",
			Message: "Slug parameter is required",
		})
		return
	}

	movie, err := c.movieHomeService.GetBySlug(ctx.Request.Context(), slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{
				Status:  "error",
				Message: "Movie not found",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Status:  "error",
			Message: "Failed to get movie",
		})
		return
	}

	// Map domain model -> new MovieDetailResponse DTO
	movieDetail := dto.MovieDetailResponse{
		ID:               movie.ID,
		Title:            movie.Name,
		ReleaseDate:      movie.ReleaseDate.Format("2006-01-02"),
		DurationInMinute: movie.DurationInMinute,
		DirectorName:     movie.DirectorName,
		Synopsis:         movie.Synopsis,
		ImagePoster:      movie.Image,
		GenresCategories: movie.Categories,
		Casts:            movie.Casts,
	}

	// Return wrapped response using string status fields
	ctx.JSON(http.StatusOK, dto.MovieDetailWrappedResponse{
		Status: "success",
		Data:   movieDetail,
	})
}
