package controller

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/tickitz-backend/internal/dto"
	"github.com/tickitz-backend/internal/jwttoken"
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

// CreatePendingOrder godoc
//
//	@Summary		Create pending order
//	@Description	Create a pending order for one exact movie cinema schedule.
//	@Tags			Orders
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.CreatePendingOrderRequest	true	"Pending order payload"
//	@Success		201		{object}	dto.SuccessResponse			"Order created successfully"
//	@Failure		400		{object}	dto.ErrorResponse			"Invalid request payload"
//	@Failure		401		{object}	dto.ErrorResponse			"Unauthorized"
//	@Failure		404		{object}	dto.ErrorResponse			"Movie cinema schedule not found"
//	@Failure		500		{object}	dto.ErrorResponse			"Failed to create order"
//	@Router			/orders [post]
func (c *orderController) CreatePendingOrder(ctx *gin.Context) {
	claims, ok := jwttoken.GetClaims(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Message: "unauthorized",
		})
		return
	}

	var req dto.CreatePendingOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "invalid request payload",
		})
		return
	}

	order, err := c.orderService.CreatePendingOrder(
		ctx.Request.Context(),
		int64(claims.UserId),
		req.MovieCinemaID,
	)
	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, dto.ErrorResponse{
				Success: false,
				Message: "movie cinema schedule not found",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Success: false,
			Message: "failed to create order",
		})
		return
	}

	statusCode := http.StatusCreated
	message := "order created successfully"
	if order.Reused {
		statusCode = http.StatusOK
		message = "pending order reused successfully"
	}

	ctx.JSON(statusCode, dto.SuccessResponse{
		Success: true,
		Message: message,
		Data:    order,
	})
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
	claims, ok := jwttoken.GetClaims(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Message: "unauthorized",
		})
		return
	}

	var req dto.OrderHistoryRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "invalid query parameters",
		})
		return
	}

	result, err := c.orderService.GetOrderHistory(
		ctx.Request.Context(),
		int64(claims.UserId),
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
