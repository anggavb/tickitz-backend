package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tickitz-backend/internal/controller"
	"github.com/tickitz-backend/internal/repository"
	"github.com/tickitz-backend/internal/service"
)

func RegisterMovieRouter(router *gin.Engine, db *pgxpool.Pool) {
	adminRouter := router.Group("/admin")
	movieRouter := adminRouter.Group("/movies")
	movieRepo := repository.NewMovieRepository(db)
	movieService := service.NewMovieService(movieRepo)
	movieController := controller.NewMovieController(movieService)

	movieRouter.GET("", movieController.List)
	movieRouter.GET("months", movieController.ListReleaseMonths)
	movieRouter.GET(":id", movieController.GetByID)
	movieRouter.POST("", movieController.Create)
	movieRouter.PATCH(":id", movieController.Update)
	movieRouter.DELETE(":id", movieController.Delete)

	adminRouter.GET("/categories", movieController.ListCategories)
	adminRouter.GET("/casts", movieController.ListCasts)
}
