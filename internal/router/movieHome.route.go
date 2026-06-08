package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tickitz-backend/internal/controller"
	"github.com/tickitz-backend/internal/repository"
	"github.com/tickitz-backend/internal/service"
)

func HomeMovieRouter(router *gin.Engine, db *pgxpool.Pool) {
	movie := router.Group("/movies")

	movieHomeRepo := repository.NewMovieHomeRepository(db)
	movieHomeService := service.NewMovieHomeService(movieHomeRepo)
	movieHomeController := controller.NewMovieHomeController(movieHomeService)

	movie.GET("/:slug", movieHomeController.GetBySlug)
	movie.GET("/:slug/schedule", movieHomeController.GetMovieSchedule)
}
