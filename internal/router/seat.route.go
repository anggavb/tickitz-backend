package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tickitz-backend/internal/controller"
	"github.com/tickitz-backend/internal/repository"
	"github.com/tickitz-backend/internal/service"
)

func RegisterSeatRouter(router *gin.Engine, db *pgxpool.Pool) {
	seatRepo := repository.NewSeatRepository(db)
	seatService := service.NewSeatService(seatRepo)
	seatController := controller.NewSeatController(seatService)

	router.GET("/movie-cinemas/:movie_cinema_id/seats", seatController.GetSeatMap)
}
