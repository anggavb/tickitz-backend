package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tickitz-backend/internal/controller"
	"github.com/tickitz-backend/internal/repository"
	"github.com/tickitz-backend/internal/service"
)

// , authCache *repository.AuthCacheRepository
func RegisterOrderRouter(router *gin.Engine, db *pgxpool.Pool) {
	orderRoute := router.Group("/orders") //  middleware.VerifyToken(authCache)
	orderRepo := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepo)
	orderController := controller.NewOrderController(orderService)

	orderRoute.GET("/history", orderController.GetOrderByUserID)
}
