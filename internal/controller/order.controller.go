package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tickitz-backend/internal/dto"
	"github.com/tickitz-backend/internal/service"
)

type orderController struct {
	orderService *service.OrderService
}

func NewOrderController(orderService *service.OrderService) *orderController {
	return &orderController{
		orderService: orderService,
	}
}

// GetOrderByUserID godoc
//
//	@Summary		Get order history by user
//	@Description	Get user order history with pagination.
//	@Tags			Orders
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int	false	"Page Number"	default(1)
//	@Param			limit	query		int	false	"Items Per Page"	default(10)
//	@Success		200		{object}	dto.SuccessResponse	"Order history retrieved successfully"
//	@Failure		400		{object}	dto.ErrorResponse		"Invalid query parameters"
//	@Failure		500		{object}	dto.ErrorResponse		"Failed to get order history"
//	@Router			/orders/history [get]
func (c *orderController) GetOrderByUserID(ctx *gin.Context) {
	// token, _ := ctx.Get("claims")
	// claims := token.(pkg.Claims)
	// userID := claims.ID
	var req dto.OrderHistoryRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "invalid query parameters",
		})
		return
	}

	// sementara dummy
	userID := int64(1)

	result, err := c.orderService.GetOrderHistory(
		ctx.Request.Context(),
		userID,
		req.Page,
		req.Limit,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "failed to get order history",
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "order history retrieved successfully",
		Data:    result,
	})
}
