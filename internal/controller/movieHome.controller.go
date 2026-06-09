package controller

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/tickitz-backend/internal/response"
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

// GetMovieBySlug godoc
// @Summary      Get movie detail and schedules by slug
// @Description  Retrieve movie detail and schedules using slug path parameter with optional date and location query parameters.
// @Tags         Movies
// @Accept       json
// @Produce      json
// @Param        slug      path      string  true   "Movie Slug"
// @Param        date      query     string  false  "Selected date in YYYY-MM-DD format" example(2026-06-09)
// @Param        location  query     string  false  "Filter by location name" example(Jakarta)
// @Success      200       {object}  dto.MovieHomeDetailWrappedResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure      500       {object}  map[string]interface{}
// @Router       /movies/{slug} [get]
func (c *MovieHomeController) GetMovieBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")
	selectedDate := ctx.Query("date")
	location := ctx.DefaultQuery("location", "")

	if slug == "" {
		response.Error(ctx, http.StatusBadRequest, "movie slug is required")
		return
	}

	if selectedDate != "" {
		if _, err := time.Parse("2006-01-02", selectedDate); err != nil {
			response.Error(ctx, http.StatusBadRequest, "date format must be YYYY-MM-DD")
			return
		}
	}

	movie, err := c.movieHomeService.GetMovieBySlug(
		ctx.Request.Context(),
		slug,
		selectedDate,
		location,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			response.Error(ctx, http.StatusNotFound, "movie not found")
			return
		}

		response.Error(ctx, http.StatusInternalServerError, "failed to get movie")
		return
	}

	response.Success(
		ctx,
		http.StatusOK,
		"movie retrieved successfully",
		movie,
	)
}
