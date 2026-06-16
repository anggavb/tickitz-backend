package controller

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/tickitz-backend/internal/dto"
	"github.com/tickitz-backend/internal/errs"
	"github.com/tickitz-backend/internal/response"
	"github.com/tickitz-backend/internal/service"
)

type MovieHomeController struct {
	movieHomeService *service.MovieHomeService
}

func NewMovieHomeController(movieHomeService *service.MovieHomeService) *MovieHomeController {
	return &MovieHomeController{movieHomeService: movieHomeService}
}

// GetMovieBySlug godoc
//
//	@Summary		Get movie detail by slug
//	@Description	Get movie detail information using movie slug.
//	@Tags			Movies
//	@Accept			json
//	@Produce		json
//	@Param			slug	path		string	true	"Movie Slug"
//	@Success		200		{object}	dto.MovieDetailSuccessResponse	"Movie detail retrieved successfully"
//	@Failure		400		{object}	dto.ErrorResponse		"Movie slug is required"
//	@Failure		404		{object}	dto.ErrorResponse		"Movie not found"
//	@Failure		500		{object}	dto.ErrorResponse		"Failed to get movie detail"
//	@Router			/movies/{slug} [get]
func (c *MovieHomeController) GetMovieBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")

	if slug == "" {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "movie slug is required",
		})
		return
	}

	movie, err := c.movieHomeService.GetMovieBySlug(ctx.Request.Context(), slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{
				Success: false,
				Message: "movie not found",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "failed to get movie detail",
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "movie detail retrieved successfully",
		Data:    movie,
	})
}

// GetMovieSchedulesBySlug godoc
//
//	@Summary		Get movie schedules by slug
//	@Description	Get available cinema schedules for a movie using movie slug.
//	@Tags			Movies
//	@Accept			json
//	@Produce		json
//	@Param			slug	path		string	true	"Movie Slug"
//	@Param			date	query		string	false	"Show Date YYYY-MM-DD"
//	@Param			time	query		string	false	"Showtime HH:MM"
//	@Param			location	query		string	false	"Location Name"
//	@Success		200		{object}	dto.MovieScheduleListSuccessResponse	"Movie schedules retrieved successfully"
//	@Failure		400		{object}	dto.ErrorResponse		"Movie slug is required"
//	@Failure		404		{object}	dto.ErrorResponse		"Movie not found"
//	@Failure		500		{object}	dto.ErrorResponse		"Failed to get movie schedules"
//	@Router			/movies/{slug}/schedules [get]
func (c *MovieHomeController) GetMovieSchedulesBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")

	if slug == "" {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "movie slug is required",
		})
		return
	}

	var filter dto.MovieScheduleQuery
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "invalid schedule filter",
		})
		return
	}

	schedules, err := c.movieHomeService.GetMovieSchedulesBySlug(ctx.Request.Context(), slug, filter)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidScheduleFilter) {
			ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Success: false,
				Message: "schedule date or time has passed",
			})
			return
		}

		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{
				Success: false,
				Message: "movie not found",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "failed to get movie schedules",
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "movie schedules retrieved successfully",
		Data:    schedules,
	})
}

// GetMovieScheduleOptionsBySlug godoc
//
//	@Summary		Get movie schedule filter options by slug
//	@Description	Get available dates, showtimes, and locations for a movie using movie slug.
//	@Tags			Movies
//	@Accept			json
//	@Produce		json
//	@Param			slug	path		string	true	"Movie Slug"
//	@Success		200		{object}	dto.MovieScheduleOptionsSuccessResponse	"Movie schedule options retrieved successfully"
//	@Failure		400		{object}	dto.ErrorResponse		"Movie slug is required"
//	@Failure		404		{object}	dto.ErrorResponse		"Movie not found"
//	@Failure		500		{object}	dto.ErrorResponse		"Failed to get movie schedule options"
//	@Router			/movies/{slug}/schedule-options [get]
func (c *MovieHomeController) GetMovieScheduleOptionsBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")

	if slug == "" {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "movie slug is required",
		})
		return
	}

	options, err := c.movieHomeService.GetMovieScheduleOptionsBySlug(ctx.Request.Context(), slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{
				Success: false,
				Message: "movie not found",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "failed to get movie schedule options",
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "movie schedule options retrieved successfully",
		Data:    options,
	})
}

// GetLocations godoc
// @Summary      Get schedule locations
// @Description  Retrieve all available movie schedule locations.
// @Tags         Movie Schedules
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.MovieLocationListSuccessResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /movies/locations [get]
func (c *MovieHomeController) GetLocations(ctx *gin.Context) {
	locations, err := c.movieHomeService.GetLocations(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to get movie locations")
		return
	}

	response.Success(
		ctx,
		http.StatusOK,
		"movie locations retrieved successfully",
		locations,
	)
}

// GetShowtimes godoc
// @Summary      Get schedule showtimes
// @Description  Retrieve all available movie schedule showtimes.
// @Tags         Movie Schedules
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.MovieShowtimeListSuccessResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /movies/showtimes [get]
func (c *MovieHomeController) GetShowtimes(ctx *gin.Context) {
	showtimes, err := c.movieHomeService.GetShowtimes(ctx.Request.Context())
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to get movie showtimes")
		return
	}

	response.Success(
		ctx,
		http.StatusOK,
		"movie showtimes retrieved successfully",
		showtimes,
	)
}

// GetMoviesWithFilter godoc
//
//	@Summary		Get movies with filter and pagination
//	@Description	Get list of movies with optional filters by category, name, and pagination.
//	@Tags			Movies
//	@Accept			json
//	@Produce		json
//	@Param	category	query		[]string	false	"Movie Categories"
//	@Param	name		query		string		false	"Movie Name"
//	@Param	showToday	query		bool		false	"Show only movies playing today"
//	@Param	page		query		int			false	"Page Number"
//	@Param	limit		query		int			false	"Items Per Page"
//	@Success		200			{object}	dto.MovieHomeListSuccessResponse
//	@Failure		400			{object}	dto.ErrorResponse
//	@Failure		500			{object}	dto.ErrorResponse
//	@Router			/movies [get]
func (c *MovieHomeController) GetMoviesWithFilter(ctx *gin.Context) {
	var param dto.MovieParamsRequest

	if err := ctx.ShouldBindQuery(&param); err != nil {
		log.Printf(
			"[MovieHomeController][GetMoviesWithFilter] bind query error: %v",
			err,
		)
		response.Error(
			ctx,
			http.StatusBadRequest,
			err.Error(),
		)
		return
	}

	data, err := c.movieHomeService.GetAllMovies(
		ctx.Request.Context(),
		param,
	)
	if err != nil {
		log.Printf(
			"[MovieHomeController][GetMoviesWithFilter] service error: %v",
			err,
		)

		response.Error(
			ctx,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	response.Success(
		ctx,
		http.StatusOK,
		"success to get movies",
		data,
	)
}

// GetUpcomingMovies godoc
//
//	@Summary		Get upcoming movies
//	@Description	Get list of upcoming movies (release date greater than current date).
//	@Tags			Movies
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dto.UpcomingMoviesSuccessResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/movies/upcoming [get]
func (c *MovieHomeController) GetUpcomingMovies(ctx *gin.Context) {

	data, err := c.movieHomeService.GetUpcomingMovies(
		ctx.Request.Context(),
	)
	if err != nil {
		log.Printf(
			"[MovieHomeController][GetUpcomingMovies] service error: %v",
			err,
		)

		response.Error(
			ctx,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	response.Success(
		ctx,
		http.StatusOK,
		"success to get upcoming movies",
		data,
	)
}
