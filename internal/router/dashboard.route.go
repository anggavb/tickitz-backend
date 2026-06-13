package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tickitz-backend/internal/controller"
	"github.com/tickitz-backend/internal/repository"
	"github.com/tickitz-backend/internal/service"
)

func RegisterDashboardRouter(router *gin.Engine, db *pgxpool.Pool) {
	adminRouter := router.Group("/admin")
	dashboardRepo := repository.NewDashboardRepository(db)
	dashboardService := service.NewDashboardService(dashboardRepo)
	dashboardController := controller.NewDashboardController(dashboardService)

	dashboardRouter := adminRouter.Group("/dashboard")
	dashboardRouter.GET("/sales-chart", dashboardController.SalesChart)
	dashboardRouter.GET("/ticket-sales", dashboardController.TicketSales)
}
