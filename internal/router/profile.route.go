package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tickitz-backend/internal/controller"
	"github.com/tickitz-backend/internal/repository"
	"github.com/tickitz-backend/internal/service"
)

func RegisterProfileRouter(router *gin.Engine, db *pgxpool.Pool) {
	profileRoute := router.Group("/profile")
	profileRepo := repository.NewProfileRepository(db)
	profileService := service.NewProfileService(profileRepo)
	profileController := controller.NewProfileController(profileService)

	profileRoute.GET("", profileController.GetProfileById)
}
