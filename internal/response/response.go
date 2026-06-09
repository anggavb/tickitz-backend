package response

import "github.com/gin-gonic/gin"

func Success(ctx *gin.Context, code int, message string, data any) {
	ctx.JSON(code, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func Error(ctx *gin.Context, code int, message string) {
	ctx.JSON(code, gin.H{
		"success": false,
		"message": message,
	})
}

type SuccessResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"success"`
	Data    any    `json:"data"`
}

type ErrorResponse struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"date format must be YYYY-MM-DD"`
}
