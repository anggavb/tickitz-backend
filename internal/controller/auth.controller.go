package controller

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tickitz-backend/internal/dto"
	"github.com/tickitz-backend/internal/errs"
	"github.com/tickitz-backend/internal/response"
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

		if strings.Contains(err.Error(), "Email") &&
			strings.Contains(err.Error(), "email") {
			response.Error(
				ctx,
				http.StatusBadRequest,
				"Email format is invalid",
			)
			return
		}

		if strings.Contains(err.Error(), "Password") &&
			strings.Contains(err.Error(), "min") {
			response.Error(
				ctx,
				http.StatusBadRequest,
				"Password must be at least 8 characters",
			)
			return
		}

		if strings.Contains(err.Error(), "IsAgree") &&
			(strings.Contains(err.Error(), "eq") ||
				strings.Contains(err.Error(), "required")) {
			response.Error(
				ctx,
				http.StatusBadRequest,
				"You must agree to the terms and conditions",
			)
			return
		}

		response.Error(
			ctx,
			http.StatusBadRequest,
			"Invalid request",
		)
		return
	}

	if err := c.authService.Register(ctx.Request.Context(), user); err != nil {

		if errors.Is(err, errs.ErrExistingEmail) {
			response.Error(
				ctx,
				http.StatusConflict,
				"Email already registered",
			)
			return
		}

		if errors.Is(err, errs.ErrAccountNotActive) {

			otpReq := dto.NewOTPRequest{
				Email: user.Email,
			}

			if otpErr := c.authService.GetNewOTP(
				ctx.Request.Context(),
				otpReq,
			); otpErr != nil {

				response.Error(
					ctx,
					http.StatusInternalServerError,
					"Failed to send new OTP",
				)
				return
			}

			response.Success(
				ctx,
				http.StatusOK,
				"Your email has been registered but not active, check your email to activate!",
				gin.H{
					"email": user.Email,
				},
			)
			return
		}

		response.Error(
			ctx,
			http.StatusInternalServerError,
			"Internal server error",
		)
		return
	}

	response.Success(
		ctx,
		http.StatusCreated,
		"Registration successful. Please check your email to activate your account.",
		gin.H{
			"email": user.Email,
		},
	)
}

func (c *AuthController) Activate(ctx *gin.Context) {
	var user dto.ActivationRequest

	if err := ctx.ShouldBindJSON(&user); err != nil {

		if strings.Contains(err.Error(), "Email") &&
			strings.Contains(err.Error(), "email") {
			response.Error(
				ctx,
				http.StatusBadRequest,
				"Email format is invalid",
			)
			return
		}

		if strings.Contains(err.Error(), "OTP") &&
			strings.Contains(err.Error(), "required") {
			response.Error(
				ctx,
				http.StatusBadRequest,
				"OTP is required",
			)
			return
		}

		response.Error(
			ctx,
			http.StatusBadRequest,
			"Invalid request",
		)
		return
	}

	if err := c.authService.Activate(
		ctx.Request.Context(),
		user,
	); err != nil {

		switch {

		case errors.Is(err, errs.ErrAccountActivated):
			response.Error(
				ctx,
				http.StatusBadRequest,
				err.Error(),
			)

		case errors.Is(err, errs.ErrTokenExpired):
			response.Error(
				ctx,
				http.StatusBadRequest,
				err.Error(),
			)

		case errors.Is(err, errs.ErrInvalidOTP):
			response.Error(
				ctx,
				http.StatusBadRequest,
				err.Error(),
			)

		case errors.Is(err, errs.ErrInternalServer):
			response.Error(
				ctx,
				http.StatusInternalServerError,
				"Internal server error",
			)

		default:
			response.Error(
				ctx,
				http.StatusBadRequest,
				err.Error(),
			)
		}

		return
	}

	response.Success(
		ctx,
		http.StatusOK,
		"Account activated successfully",
		gin.H{
			"email": user.Email,
		},
	)
}

func (c *AuthController) GetNewOTP(ctx *gin.Context) {
	var user dto.NewOTPRequest

	if err := ctx.ShouldBindJSON(&user); err != nil {

		if strings.Contains(err.Error(), "Email") &&
			strings.Contains(err.Error(), "email") {
			response.Error(
				ctx,
				http.StatusBadRequest,
				"Email format is invalid",
			)
			return
		}

		response.Error(
			ctx,
			http.StatusBadRequest,
			"Invalid request",
		)
		return
	}

	if err := c.authService.GetNewOTP(
		ctx.Request.Context(),
		user,
	); err != nil {

		if errors.Is(err, errs.ErrAccountActivated) {
			response.Error(
				ctx,
				http.StatusConflict,
				"Account already activated",
			)
			return
		}

		response.Error(
			ctx,
			http.StatusInternalServerError,
			"Internal server error",
		)
		return
	}

	response.Success(
		ctx,
		http.StatusOK,
		"New OTP has been sent successfully",
		gin.H{
			"email": user.Email,
		},
	)
}

func (c *AuthController) Login(ctx *gin.Context) {
	var user dto.LoginRequest

	if err := ctx.ShouldBindJSON(&user); err != nil {

		if strings.Contains(err.Error(), "Email") &&
			strings.Contains(err.Error(), "email") {
			response.Error(
				ctx,
				http.StatusBadRequest,
				"Email format is invalid",
			)
			return
		}

		if strings.Contains(err.Error(), "Password") &&
			strings.Contains(err.Error(), "min") {
			response.Error(
				ctx,
				http.StatusBadRequest,
				"Password must be at least 8 characters",
			)
			return
		}

		response.Error(
			ctx,
			http.StatusBadRequest,
			"Invalid request",
		)
		return
	}

	data, err := c.authService.Login(ctx.Request.Context(), user.Email, user.Password)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(
		ctx,
		http.StatusOK,
		"Login success",
		data,
	)
}
func (c *AuthController) ChangeUserPassword(ctx *gin.Context) {
	var req dto.ChangePasswordRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("[ChangeUserPassword] BindJSON error: %v\n", err)

		if strings.Contains(err.Error(), "NewPassword") {
			response.Error(ctx, http.StatusBadRequest, err.Error())
			return
		}

		response.Error(ctx, http.StatusBadRequest, "invalid request")
		return
	}

	// TODO: ambil dari JWT claims
	userID := 1

	if err := c.authService.ChangeUserPassword(
		ctx.Request.Context(),
		req.NewPassword,
		userID,
	); err != nil {
		log.Printf("[ChangeUserPassword] Service error: %v\n", err)

		response.Error(
			ctx,
			http.StatusInternalServerError,
			"failed to change password",
		)
		return
	}

	response.Success(ctx, http.StatusOK, "password updated", nil)
}
