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
