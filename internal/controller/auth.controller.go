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

		if errors.Is(err, errs.ErrAccountNotActive) {

			otpReq := dto.NewOTPRequest{
				Email: user.Email,
			}

			if otpErr := c.authService.GetNewOTP(ctx.Request.Context(), otpReq); otpErr != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Failed to send new OTP",
				})
				return
			}

			ctx.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "Your email has been registered but not active, check your email to activate!",
				"data": gin.H{
					"email": user.Email,
				},
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

func (c *AuthController) Activate(ctx *gin.Context) {
	var user dto.ActivationRequest

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
		if strings.Contains(err.Error(), "OTP") {
			if strings.Contains(err.Error(), "required") {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "OTP is required",
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
	if err := c.authService.Activate(ctx, user); err != nil {
		switch err {
		case errs.ErrAccountActivated:
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
			})

		case errs.ErrTokenExpired:
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
			})

		case errs.ErrInvalidOTP:
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
			})

		case errs.ErrInternalServer:
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": err.Error(),
			})

		default:
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
			})
		}

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Account activated successfully",
		"data": gin.H{
			"email": user.Email,
		},
	})

}

func (c *AuthController) GetNewOTP(ctx *gin.Context) {
	var user dto.NewOTPRequest

	if err := ctx.ShouldBindJSON(&user); err != nil {

		if strings.Contains(err.Error(), "Email") &&
			strings.Contains(err.Error(), "email") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Email format is invalid",
			})
			return
		}

		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request",
		})
		return
	}

	if err := c.authService.GetNewOTP(ctx.Request.Context(), user); err != nil {

		if errors.Is(err, errs.ErrAccountActivated) {
			ctx.JSON(http.StatusConflict, gin.H{
				"success": false,
				"message": "Account already activated",
			})
			return
		}

		if errors.Is(err, errs.ErrAccountNotActive) {
			ctx.JSON(http.StatusConflict, gin.H{
				"success": false,
				"message": "Email hasn't activate yet",
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

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "New OTP has been sent successfully",
		"data": gin.H{
			"email": user.Email,
		},
	})
}
