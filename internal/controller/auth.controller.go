package controller

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tickitz-backend/internal/dto"
	"github.com/tickitz-backend/internal/errs"
	"github.com/tickitz-backend/internal/service"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (c *AuthController) Register(ctx *gin.Context) {
	var user dto.RegisterRequest

	if err := ctx.ShouldBindJSON(&user); err != nil {

		if strings.Contains(err.Error(), "Email") {
			if strings.Contains(err.Error(), "email") {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "Email format is invalid",
				})
				return
			}
		}

		if strings.Contains(err.Error(), "Password") {
			if strings.Contains(err.Error(), "min") {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "Password must be at least 8 characters",
				})
				return
			}
		}

		if strings.Contains(err.Error(), "IsAgree") {
			if strings.Contains(err.Error(), "eq") || strings.Contains(err.Error(), "required") {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "You must agree to the terms and conditions",
				})
				return
			}
		}

		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request",
		})
		return
	}

	if err := c.authService.Register(ctx.Request.Context(), user); err != nil {

		if errors.Is(err, errs.ErrExistingEmail) {
			ctx.JSON(http.StatusConflict, gin.H{
				"success": false,
				"message": "Email already registered",
			})
			return
		}

		if errors.Is(err, errs.ErrInternalServer) {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Internal server error",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal server error",
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Registration successful. Please check your email to activate your account.",
		"data": gin.H{
			"email": user.Email,
		},
	})
}
