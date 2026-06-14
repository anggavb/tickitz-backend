package controller

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/tickitz-backend/internal/dto"
	"github.com/tickitz-backend/internal/jwttoken"
	"github.com/tickitz-backend/internal/repository"
	"github.com/tickitz-backend/internal/service"
	qrcode "github.com/yeqown/go-qrcode/v2"
)

type orderController struct {
	orderService *service.OrderService
}

func NewOrderController(orderService *service.OrderService) *orderController {
	return &orderController{
		orderService: orderService,
	}
}

func getUserID(ctx *gin.Context) (int64, bool) {
	claims, ok := jwttoken.GetClaims(ctx)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Success: false,
			Message: "unauthorized",
		})
		return 0, false
	}

	return int64(claims.UserId), true
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
	userID, ok := getUserID(ctx)
	if !ok {
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
		userID,
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

func (c *orderController) GetOrderDetail(ctx *gin.Context) {
	userID, ok := getUserID(ctx)
	if !ok {
		return
	}

	order, err := c.orderService.GetOrderDetail(ctx.Request.Context(), userID, ctx.Param("order_id"))
	if err != nil {
		writeOrderError(ctx, err, "failed to get order detail")
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "order detail retrieved successfully",
		Data:    order,
	})
}

func (c *orderController) UpdateOrderSeats(ctx *gin.Context) {
	userID, ok := getUserID(ctx)
	if !ok {
		return
	}

	var req dto.UpdateOrderSeatsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "invalid request payload",
		})
		return
	}

	order, err := c.orderService.UpdateOrderSeats(
		ctx.Request.Context(),
		userID,
		ctx.Param("order_id"),
		req.Seats,
	)
	if err != nil {
		writeOrderError(ctx, err, "failed to update order seats")
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "order seats updated successfully",
		Data:    order,
	})
}

func (c *orderController) GetPaymentMethods(ctx *gin.Context) {
	userID, ok := getUserID(ctx)
	if !ok {
		return
	}

	if _, err := c.orderService.GetOrderDetail(ctx.Request.Context(), userID, ctx.Param("order_id")); err != nil {
		writeOrderError(ctx, err, "failed to get payment methods")
		return
	}

	methods, err := c.orderService.GetPaymentMethods(ctx.Request.Context())
	if err != nil {
		writeOrderError(ctx, err, "failed to get payment methods")
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "payment methods retrieved successfully",
		Data:    methods,
	})
}

func (c *orderController) SubmitPayment(ctx *gin.Context) {
	userID, ok := getUserID(ctx)
	if !ok {
		return
	}

	var req dto.UpdateOrderPaymentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Success: false,
			Message: "invalid request payload",
		})
		return
	}

	order, err := c.orderService.SubmitPayment(
		ctx.Request.Context(),
		userID,
		ctx.Param("order_id"),
		req,
	)
	if err != nil {
		writeOrderError(ctx, err, "failed to submit payment")
		return
	}

	ctx.JSON(http.StatusOK, dto.SuccessResponse{
		Success: true,
		Message: "payment completed successfully",
		Data:    order,
	})
}

func (c *orderController) GetOrderQR(ctx *gin.Context) {
	userID, ok := getUserID(ctx)
	if !ok {
		return
	}

	payload, err := c.orderService.BuildTicketQRPayload(ctx.Request.Context(), userID, ctx.Param("order_id"))
	if err != nil {
		writeOrderError(ctx, err, "failed to generate ticket qr")
		return
	}

	qr, err := qrcode.New(string(payload))
	if err != nil {
		writeOrderError(ctx, err, "failed to generate ticket qr")
		return
	}

	var buf bytes.Buffer
	if err := qr.Save(newPNGQRWriter(&buf)); err != nil {
		writeOrderError(ctx, err, "failed to generate ticket qr")
		return
	}

	ctx.Header("Cache-Control", "no-store")
	ctx.Data(http.StatusOK, "image/png", buf.Bytes())
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
	userID, ok := getUserID(ctx)
	if !ok {
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

func writeOrderError(ctx *gin.Context, err error, fallbackMessage string) {
	switch {
	case errors.Is(err, repository.ErrOrderNotFound), errors.Is(err, pgx.ErrNoRows):
		ctx.JSON(http.StatusNotFound, dto.ErrorResponse{Success: false, Message: "order not found"})
	case errors.Is(err, repository.ErrInvalidSeats):
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Message: "invalid seats"})
	case errors.Is(err, repository.ErrSeatUnavailable):
		ctx.JSON(http.StatusConflict, dto.ErrorResponse{Success: false, Message: "some selected seats are no longer available"})
	case errors.Is(err, repository.ErrPaymentMethodNotFound):
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Success: false, Message: "payment method not found"})
	case errors.Is(err, repository.ErrOrderNotPayable):
		ctx.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{Success: false, Message: "order is not ready for payment"})
	case errors.Is(err, repository.ErrOrderAlreadyPaid):
		ctx.JSON(http.StatusUnprocessableEntity, dto.ErrorResponse{Success: false, Message: "order already paid"})
	default:
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Success: false, Message: fallbackMessage})
	}
}

type pngQRWriter struct {
	buf    *bytes.Buffer
	scale  int
	margin int
}

func newPNGQRWriter(buf *bytes.Buffer) *pngQRWriter {
	return &pngQRWriter{
		buf:    buf,
		scale:  10,
		margin: 4,
	}
}

func (w *pngQRWriter) Close() error {
	return nil
}

func (w *pngQRWriter) Write(mat qrcode.Matrix) error {
	bitmap := mat.Bitmap()
	moduleCount := mat.Width()
	imageSize := (moduleCount + w.margin*2) * w.scale
	img := image.NewRGBA(image.Rect(0, 0, imageSize, imageSize))
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	black := color.RGBA{A: 255}

	for y := 0; y < imageSize; y++ {
		for x := 0; x < imageSize; x++ {
			img.Set(x, y, white)
		}
	}

	for y, row := range bitmap {
		for x, set := range row {
			if !set {
				continue
			}
			startX := (x + w.margin) * w.scale
			startY := (y + w.margin) * w.scale
			for py := 0; py < w.scale; py++ {
				for px := 0; px < w.scale; px++ {
					img.Set(startX+px, startY+py, black)
				}
			}
		}
	}

	return png.Encode(w.buf, img)
}
