package controller

import (
	"database/sql"
	"errors"
	"fmt"
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
// @Summary      Get a movie by its slug
// @Description  Retrieve detailed information about a specific movie using its unique URI slug.
// @Tags         movies
// @Accept       json
// @Produce      json
// @Param slug path string true "Movie Slug"
// @Success 200 {object} dto.MovieSingleResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router       /movies/{slug} [get]
func (c *MovieHomeController) GetBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")

	movie, err := c.movieHomeService.GetBySlug(ctx.Request.Context(), slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Movie not found",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get movie",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    buildMovieResponse(movie),
	})
}

// GetMovieSchedule handles retrieving movie showtimes by slug and location.
//
//	@Summary		Get movie showtimes schedule
//	@Description	Retrieve nested showtimes for a specific movie grouped by locations and cinema branches.
//	@Tags			Movies
//	@Accept			json
//	@Produce		json
//	@Param			slug		path		string	true	"Movie Slug (e.g., echoes-of-jakarta)"
//	@Param			location	query		string	false	"Filter by specific city/region location name"
//	@Success		200			{object}	dto.MovieScheduleWrappedResponse	"Successfully retrieved schedules"
//	@Failure		404			{object}	dto.ErrorResponse					"Movie or location not found"
//	@Failure		500			{object}	dto.ErrorResponse					"Internal server error"
//	@Router			/movies/{slug}/schedule [get]
func (ctrl *MovieHomeController) GetMovieSchedule(c *gin.Context) {
	slug := c.Param("slug")
	location := c.Query("location")

	schedules, err := ctrl.movieHomeService.GetScheduleBySlugAndLocation(c.Request.Context(), slug, location)
	if err != nil {
		// If using standard database/sql, check for sql.ErrNoRows.
		// Alternatively, check against your own domain error (e.g., service.ErrMovieNotFound).
		if errors.Is(err, sql.ErrNoRows) || err.Error() == "movie not found" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Status:  "error",
				Message: fmt.Sprintf("Movie schedule for slug '%s' could not be found", slug),
			})
			return
		}

		// Fallback for true system failures
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Status:  "error",
			Message: "Internal server processing error",
		})
		return
	}

	c.JSON(http.StatusOK, dto.MovieScheduleWrappedResponse{
		Status: "success",
		Data:   schedules,
	})
}
