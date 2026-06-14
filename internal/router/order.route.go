package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tickitz-backend/internal/controller"
	"github.com/tickitz-backend/internal/middleware"
	"github.com/tickitz-backend/internal/repository"
	"github.com/tickitz-backend/internal/service"
)

func RegisterOrderRouter(router *gin.Engine, db *pgxpool.Pool, authCache *repository.AuthCacheRepository) {
	orderRoute := router.Group("/orders", middleware.VerifyToken(authCache))
	orderRepo := repository.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepo)
	orderController := controller.NewOrderController(orderService)

	orderRoute.POST("", orderController.CreatePendingOrder)
	orderRoute.GET("/history", orderController.GetOrderByUserID)
	orderRoute.GET("/:order_id", orderController.GetOrderDetail)
	orderRoute.PATCH("/:order_id/seats", orderController.UpdateOrderSeats)
	orderRoute.GET("/:order_id/payment-methods", orderController.GetPaymentMethods)
	orderRoute.PATCH("/:order_id/payment", orderController.SubmitPayment)
	orderRoute.GET("/:order_id/qr", orderController.GetOrderQR)
}
