package controller

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
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
// @Param        slug      path      string  true  "Movie Slug"
// @Success      200       {object}  map[string]interface{} "Success response with movie data"
// @Failure      404       {object}  map[string]interface{} "Movie not found"
// @Failure      500       {object}  map[string]interface{} "Internal server error"
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
