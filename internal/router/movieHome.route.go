package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tickitz-backend/internal/controller"
	"github.com/tickitz-backend/internal/repository"
	"github.com/tickitz-backend/internal/service"
)

func HomeMovieRouter(router *gin.Engine, db *pgxpool.Pool) {
	movieGroup := router.Group("/movies")

	movieHomeRepository := repository.NewMovieHomeRepository(db)
	movieScheduleRepository := repository.NewMovieScheduleRepository(db)

	movieHomeService := service.NewMovieHomeService(movieHomeRepository, movieScheduleRepository)
	movieHomeController := controller.NewMovieHomeController(movieHomeService)

	movieGroup.GET("/:slug", movieHomeController.GetMovieBySlug)
	movieGroup.GET("/:slug/schedules", movieHomeController.GetMovieSchedulesBySlug)
}
