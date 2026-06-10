package response

import (
	"github.com/gin-gonic/gin"
	"github.com/tickitz-backend/internal/dto"
)

func Success(ctx *gin.Context, code int, message string, data any) {
	ctx.JSON(code, dto.SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(ctx *gin.Context, code int, message string) {
	ctx.JSON(code, dto.ErrorResponse{
		Success: false,
		Message: message,
	})
}
