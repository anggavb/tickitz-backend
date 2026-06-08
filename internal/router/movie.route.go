package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tickitz-backend/internal/controller"
	"github.com/tickitz-backend/internal/repository"
	"github.com/tickitz-backend/internal/service"
)

func RegisterMovieRouter(router *gin.Engine, db *pgxpool.Pool) {
	movieRouter := router.Group("/admin/movies")
	movieRepo := repository.NewMovieRepository(db)
	movieService := service.NewMovieService(movieRepo)
	movieController := controller.NewMovieController(movieService)

	movieRouter.GET("", movieController.List)
	movieRouter.GET(":id", movieController.GetByID)
	movieRouter.POST("", movieController.Create)
	movieRouter.PUT(":id", movieController.Update)
	movieRouter.DELETE(":id", movieController.Delete)
}
