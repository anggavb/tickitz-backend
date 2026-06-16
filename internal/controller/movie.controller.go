package controller

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tickitz-backend/internal/dto"
	"github.com/tickitz-backend/internal/model"
	"github.com/tickitz-backend/internal/service"
)

type MovieController struct {
	movieService *service.MovieService
}

func NewMovieController(movieService *service.MovieService) *MovieController {
	return &MovieController{movieService: movieService}
}

// ListMovies godoc
// @Summary Get list of movies
// @Description Get paginated list of movies
// @Tags Admin Movies
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Param month query string false "Release month filter (YYYY-MM)"
// @Security ApiKeyAuth
// @Success 200 {object} dto.MovieListResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /admin/movies [get]
func (c *MovieController) List(ctx *gin.Context) {
	page := 1
	limit := 5
	releaseMonth := ctx.Query("month")

	if pageParam := ctx.Query("page"); pageParam != "" {
		if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
			page = p
		}
	}

	if limitParam := ctx.Query("limit"); limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 {
			limit = l
		}
	}

	movies, pagination, err := c.movieService.List(ctx.Request.Context(), page, limit, releaseMonth)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to list movies",
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.MovieListResponse{
		Success: true,
		Data:    buildMovieListResponse(movies),
		Pagination: dto.Meta{
			Page:      pagination.Page,
			Limit:     pagination.Limit,
			TotalData: pagination.TotalData,
			TotalPage: pagination.TotalPage,
		},
	})
}

// GetMovieByID godoc
// @Summary Get movie by ID
// @Description Get movie by ID
// @Tags Admin Movies
// @Accept json
// @Produce json
// @Param id path int true "Movie ID"
// @Security ApiKeyAuth
// @Success 200 {object} dto.MovieSingleResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /admin/movies/{id} [get]
func (c *MovieController) GetByID(ctx *gin.Context) {
	movieID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil || movieID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid movie ID",
		})
		return
	}

	movie, err := c.movieService.GetByID(ctx.Request.Context(), movieID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Movie not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    buildMovieResponse(movie),
	})
}

// CreateMovie godoc
//
//	@Summary		Create movie
//	@Description	Create a movie. Multipart form supports optional image upload; JSON with the same non-file fields is also accepted.
//	@Tags			Admin Movies
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			name				formData	string	true	"Movie name"
//	@Param			release_date		formData	string	true	"Release date (YYYY-MM-DD)"
//	@Param			duration_in_minute	formData	int		true	"Duration in minutes"
//	@Param			director_name		formData	string	false	"Director name"
//	@Param			synopsis			formData	string	false	"Synopsis"
//	@Param			categories			formData	[]string	false	"Categories"
//	@Param			cast				formData	[]string	false	"Casts"
//	@Param			image				formData	file	false	"Movie poster image"
//	@Success		201					{object}	dto.MovieSingleResponse
//	@Failure		400					{object}	dto.ErrorResponse
//	@Failure		401					{object}	dto.ErrorResponse
//	@Failure		403					{object}	dto.ErrorResponse
//	@Router			/admin/movies [post]
func (c *MovieController) Create(ctx *gin.Context) {
	request, err := c.bindMovieRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	movie, err := c.movieService.Create(ctx.Request.Context(), request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid movie data",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    buildMovieResponse(movie),
	})
}

// ListReleaseMonths godoc
//
//	@Summary		List movie release months
//	@Description	Get available release months from movies.
//	@Tags			Admin Movies
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	dto.StringListSuccessResponse
//	@Failure		401	{object}	dto.ErrorResponse
//	@Failure		403	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/admin/movies/months [get]
func (c *MovieController) ListReleaseMonths(ctx *gin.Context) {
	months, err := c.movieService.ListReleaseMonths(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch release months",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    months,
	})
}

// ListCategories godoc
//
//	@Summary		List categories
//	@Description	Get all movie categories for admin forms.
//	@Tags			Admin Movies
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	dto.StringListSuccessResponse
//	@Failure		401	{object}	dto.ErrorResponse
//	@Failure		403	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/admin/categories [get]
func (c *MovieController) ListCategories(ctx *gin.Context) {
	categories, err := c.movieService.ListCategories(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch categories",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    categories,
	})
}

// ListCasts godoc
//
//	@Summary		List casts
//	@Description	Get all movie casts for admin forms.
//	@Tags			Admin Movies
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	dto.StringListSuccessResponse
//	@Failure		401	{object}	dto.ErrorResponse
//	@Failure		403	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/admin/casts [get]
func (c *MovieController) ListCasts(ctx *gin.Context) {
	casts, err := c.movieService.ListCasts(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch casts",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    casts,
	})
}

// GetCinemas godoc
//
//	@Summary		List cinemas
//	@Description	Get all cinemas for admin showtime forms.
//	@Tags			Admin Movies
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	dto.CinemaListSuccessResponse
//	@Failure		401	{object}	dto.ErrorResponse
//	@Failure		403	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/admin/cinemas [get]
func (c *MovieController) GetCinemas(ctx *gin.Context) {
	cinemas, err := c.movieService.GetAllCinemas(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch cinemas",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    cinemas,
	})
}

// GetMovieShowtimes godoc
//
//	@Summary		Get movie showtimes
//	@Description	Get configured cinema schedules for a movie by movie ID.
//	@Tags			Admin Movies
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		int	true	"Movie ID"
//	@Success		200	{object}	dto.MovieScheduleListSuccessResponse
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		401	{object}	dto.ErrorResponse
//	@Failure		403	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/admin/movies/{id}/showtimes [get]
func (c *MovieController) GetMovieShowtimes(ctx *gin.Context) {
	movieID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil || movieID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid movie ID"})
		return
	}

	schedules, err := c.movieService.GetMovieShowtimes(ctx.Request.Context(), movieID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to fetch showtimes"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "data": schedules})
}

// AddMovieShowtimes godoc
//
//	@Summary		Add movie showtimes
//	@Description	Create or update cinema schedules for a movie.
//	@Tags			Admin Movies
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id		path	int								true	"Movie ID"
//	@Param			payload	body	dto.AdminMovieShowtimesRequest	true	"Movie showtimes payload"
//	@Success		201		{object}	dto.EmptyDataResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		401		{object}	dto.ErrorResponse
//	@Failure		403		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/admin/movies/{id}/showtimes [post]
func (c *MovieController) AddMovieShowtimes(ctx *gin.Context) {
	movieID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil || movieID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid movie ID"})
		return
	}

	var req dto.AdminMovieShowtimesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Invalid request payload"})
		return
	}

	if req.Price < 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Price must be non-negative"})
		return
	}

	if err := c.movieService.AddMovieShowtimes(ctx.Request.Context(), movieID, req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"success": true, "message": "Showtimes saved successfully"})
}

// UpdateMovie godoc
//
//	@Summary		Update movie
//	@Description	Update a movie. Multipart form supports optional image upload; JSON with the same non-file fields is also accepted.
//	@Tags			Admin Movies
//	@Accept			multipart/form-data
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id					path		int		true	"Movie ID"
//	@Param			name				formData	string	true	"Movie name"
//	@Param			release_date		formData	string	true	"Release date (YYYY-MM-DD)"
//	@Param			duration_in_minute	formData	int		true	"Duration in minutes"
//	@Param			director_name		formData	string	false	"Director name"
//	@Param			synopsis			formData	string	false	"Synopsis"
//	@Param			categories			formData	[]string	false	"Categories"
//	@Param			cast				formData	[]string	false	"Casts"
//	@Param			image				formData	file	false	"Movie poster image"
//	@Success		200					{object}	dto.MovieSingleResponse
//	@Failure		400					{object}	dto.ErrorResponse
//	@Failure		401					{object}	dto.ErrorResponse
//	@Failure		403					{object}	dto.ErrorResponse
//	@Failure		404					{object}	dto.ErrorResponse
//	@Router			/admin/movies/{id} [patch]
func (c *MovieController) Update(ctx *gin.Context) {
	movieID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil || movieID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid movie ID",
		})
		return
	}

	existing, err := c.movieService.GetByID(ctx.Request.Context(), movieID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Movie not found",
		})
		return
	}

	request, err := c.bindMovieRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if request.Image == "" {
		request.Image = existing.Image
	}

	updatedMovie, err := c.movieService.Update(ctx.Request.Context(), movieID, request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid movie data",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    buildMovieResponse(updatedMovie),
	})
}

func (c *MovieController) bindMovieRequest(ctx *gin.Context) (dto.MovieRequest, error) {
	var request dto.MovieRequest
	contentType := ctx.ContentType()
	if strings.HasPrefix(contentType, "multipart/form-data") {
		// Parse form values manually to avoid binding errors for file fields
		if err := ctx.Request.ParseMultipartForm(32 << 20); err != nil && err != http.ErrNotMultipart {
			return request, err
		}

		// simple fields
		request.Name = ctx.PostForm("name")
		request.ReleaseDate = ctx.PostForm("release_date")
		request.DirectorName = ctx.PostForm("director_name")
		request.Synopsis = ctx.PostForm("synopsis")

		// duration: parse int
		if v := ctx.PostForm("duration_in_minute"); v != "" {
			if iv, err := strconv.Atoi(v); err == nil {
				request.DurationInMinute = iv
			} else {
				return request, err
			}
		}

		// categories: allow multiple values or comma-separated values
		request.Categories = ctx.PostFormArray("categories")
		if len(request.Categories) == 0 {
			request.Categories = parseCommaSeparatedField(ctx.PostForm("categories"))
		}

		// casts: allow multiple values or comma-separated values
		request.Casts = ctx.PostFormArray("cast")
		if len(request.Casts) == 0 {
			request.Casts = parseCommaSeparatedField(ctx.PostForm("cast"))
		}
		if len(request.Casts) == 0 {
			request.Casts = ctx.PostFormArray("casts")
		}

		// file (optional)
		file, err := ctx.FormFile("image")
		if err == nil && file != nil {
			imagePath, err := c.saveUploadedImage(ctx, file)
			if err != nil {
				return request, err
			}
			request.Image = imagePath
		}

		// Basic validation: required fields
		if request.Name == "" || request.ReleaseDate == "" || request.DurationInMinute <= 0 {
			return request, fmt.Errorf("missing required form fields")
		}

		return request, nil
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		return request, err
	}
	return request, nil
}

func parseCommaSeparatedField(value string) []string {
	if value == "" {
		return nil
	}

	parts := strings.FieldsFunc(value, func(r rune) bool {
		return r == ',' || r == ';'
	})

	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

func (c *MovieController) saveUploadedImage(ctx *gin.Context, file *multipart.FileHeader) (string, error) {
	uploadDir := filepath.Join("public", "img", "movies")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", err
	}

	ext := filepath.Ext(file.Filename)
	fileName := uuid.NewString() + ext
	destination := filepath.Join(uploadDir, fileName)

	if err := ctx.SaveUploadedFile(file, destination); err != nil {
		return "", err
	}

	return "/img/movies/" + fileName, nil
}

// DeleteMovie godoc
//
//	@Summary		Delete movie
//	@Description	Delete a movie by ID.
//	@Tags			Admin Movies
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			id	path		int	true	"Movie ID"
//	@Success		200	{object}	dto.EmptyDataResponse
//	@Failure		400	{object}	dto.ErrorResponse
//	@Failure		401	{object}	dto.ErrorResponse
//	@Failure		403	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/admin/movies/{id} [delete]
func (c *MovieController) Delete(ctx *gin.Context) {
	movieID, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil || movieID <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid movie ID",
		})
		return
	}

	if err := c.movieService.Delete(ctx.Request.Context(), movieID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete movie",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Movie deleted successfully",
	})
}

func buildMovieListResponse(movies []model.Movie) []dto.MovieResponse {
	result := make([]dto.MovieResponse, 0, len(movies))
	for _, movie := range movies {
		result = append(result, buildMovieResponse(movie))
	}
	return result
}

func buildMovieResponse(movie model.Movie) dto.MovieResponse {
	response := dto.MovieResponse{
		ID:               movie.ID,
		Name:             movie.Name,
		Slug:             movie.Slug,
		ReleaseDate:      movie.ReleaseDate.Format("2006-01-02"),
		DurationInMinute: movie.DurationInMinute,
		DirectorName:     movie.DirectorName,
		Synopsis:         movie.Synopsis,
		Image:            movie.Image,
		Categories:       movie.Categories,
		Casts:            movie.Casts,
	}

	if !movie.CreatedAt.IsZero() {
		response.CreatedAt = movie.CreatedAt.Format(time.RFC3339)
	}
	if movie.UpdatedAt != nil {
		response.UpdatedAt = movie.UpdatedAt.Format(time.RFC3339)
	}

	return response
}
