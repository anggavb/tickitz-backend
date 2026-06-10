package controller

import (
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
	return &MovieHomeController{movieHomeService: movieHomeService}
}
func (c *MovieHomeController) GetMovieBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")

	if slug == "" {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Message: "movie slug is required"})
		return
	}

	movie, err := c.movieHomeService.GetMovieBySlug(ctx.Request.Context(), slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{Success: false, Message: "movie not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Success: false, Message: "failed to get movie detail"})
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "movie detail retrieved successfully",
		Data:    movie,
	})
}

func (c *MovieHomeController) GetMovieSchedulesBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")

	if slug == "" {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Message: "movie slug is required"})
		return
	}

	schedules, err := c.movieHomeService.GetMovieSchedulesBySlug(ctx.Request.Context(), slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{Success: false, Message: "movie not found"})
			return
		}

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Success: false, Message: "failed to get movie schedules"})
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "movie schedules retrieved successfully",
		Data:    schedules,
	})
}
