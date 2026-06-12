package controller

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tickitz-backend/internal/dto"
	"github.com/tickitz-backend/internal/repository"
	"github.com/tickitz-backend/internal/service"
)

type DashboardController struct {
	dashboardService *service.DashboardService
}

func NewDashboardController(dashboardService *service.DashboardService) *DashboardController {
	return &DashboardController{dashboardService: dashboardService}
}

// SalesChart godoc
// @Summary Get admin sales chart
// @Description Get revenue chart data for the admin dashboard filtered by movie name and period
// @Tags Dashboard
// @Accept json
// @Produce json
// @Param movie_name query string false "Movie name filter"
// @Param period query string false "Period (daily, weekly, monthly)"
// @Success 200 {object} dto.DashboardSalesResponse
// @Failure 500 {object} map[string]interface{}
// @Router /admin/dashboard/sales-chart [get]
func (c *DashboardController) SalesChart(ctx *gin.Context) {
	var query dto.DashboardSalesChartRequest
	if err := ctx.BindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request query",
		})
		return
	}

	points, normalizedPeriod, err := c.dashboardService.GetSalesChart(ctx.Request.Context(), query.MovieName, query.Period)
	if err != nil {
		log.Printf("[Dashboard] SalesChart error: %v\n", err)
		if errors.Is(err, repository.ErrInvalidDashboardPeriod) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid period value; valid values are daily, weekly, monthly",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch sales chart data",
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.DashboardSalesResponse{
		Success: true,
		Data: dto.DashboardSalesData{
			MovieName: query.MovieName,
			Period:    normalizedPeriod,
			Points:    points,
		},
	})
}

// TicketSales godoc
// @Summary Get admin ticket sales data
// @Description Get ticket sales chart data for the admin dashboard filtered by category and location
// @Tags Dashboard
// @Accept json
// @Produce json
// @Param category query string false "Movie category filter"
// @Param location query string false "Cinema location filter"
// @Param period query string false "Period (daily, weekly, monthly)"
// @Success 200 {object} dto.DashboardTicketSalesResponse
// @Failure 500 {object} map[string]interface{}
// @Router /admin/dashboard/ticket-sales [get]
func (c *DashboardController) TicketSales(ctx *gin.Context) {
	var query dto.DashboardTicketSalesRequest
	if err := ctx.BindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request query",
		})
		return
	}

	points, _, err := c.dashboardService.GetTicketSales(ctx.Request.Context(), query.Category, query.Location, query.Period)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to fetch ticket sales data",
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.DashboardTicketSalesResponse{
		Success: true,
		Data:    points,
	})
}
