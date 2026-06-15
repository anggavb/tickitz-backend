package controller

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tickitz-backend/internal/dto"
	"github.com/tickitz-backend/internal/errs"
	"github.com/tickitz-backend/internal/jwttoken"
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

// Register godoc
//
//	@Summary		Register user
//	@Description	Register a new user account and send activation OTP to email.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.RegisterRequest	true	"Register payload"
//	@Success		201		{object}	dto.AuthEmailSuccessResponse
//	@Success		200		{object}	dto.AuthEmailSuccessResponse	"Email already registered but account is not active"
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		409		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/auth/signup [post]
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

// Activate godoc
//
//	@Summary		Activate user account
//	@Description	Activate a registered account using email and OTP.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.ActivationRequest	true	"Activation payload"
//	@Success		200		{object}	dto.AuthEmailSuccessResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/auth/activate [post]
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

// GetNewOTP godoc
//
//	@Summary		Request new OTP
//	@Description	Send a new activation OTP to an inactive account email.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.NewOTPRequest	true	"New OTP payload"
//	@Success		200		{object}	dto.AuthEmailSuccessResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		409		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/auth/otp [post]
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

// Login godoc
//
//	@Summary		Sign in
//	@Description	Authenticate user and return JWT token.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.LoginRequest	true	"Login payload"
//	@Success		200		{object}	dto.LoginSuccessResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Router			/auth/signin [post]
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

// Logout godoc
//
//	@Summary		Logout
//	@Description	Invalidate the current bearer token.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Success		200	{object}	dto.EmptyDataResponse
//	@Failure		401	{object}	dto.ErrorResponse
//	@Failure		500	{object}	dto.ErrorResponse
//	@Router			/auth/logout [delete]
func (c *AuthController) Logout(ctx *gin.Context) {
	claims, ok := jwttoken.GetClaims(ctx)
	if !ok {
		response.Error(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}

	tokenHashValue, ok := ctx.Get("token_hash")
	if !ok {
		response.Error(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}

	tokenHash, ok := tokenHashValue.(string)
	if !ok {
		response.Error(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := c.authService.Logout(ctx.Request.Context(), tokenHash, claims.UserId); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Internal server error")
		return
	}

	response.Success(ctx, http.StatusOK, "Logout success", nil)
}

// ChangeUserPassword godoc
//
//	@Summary		Change password
//	@Description	Change password for the currently authenticated user.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Security		ApiKeyAuth
//	@Param			payload	body		dto.ChangePasswordRequest	true	"Change password payload"
//	@Success		200		{object}	dto.EmptyDataResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		401		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/auth/password [patch]
func (c *AuthController) ChangeUserPassword(ctx *gin.Context) {
	claims, ok := jwttoken.GetClaims(ctx)
	if !ok {
		response.Error(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userID := claims.UserId

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

// ForgotPassword godoc
//
//	@Summary		Forgot password
//	@Description	Send a password reset link to the requested email.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		dto.ForgotPasswordRequest	true	"Forgot password payload"
//	@Success		200		{object}	dto.EmptyDataResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/auth/password/forgot [post]
func (c *AuthController) ForgotPassword(ctx *gin.Context) {
	var user dto.ForgotPasswordRequest
	if err := ctx.ShouldBindJSON(&user); err != nil {
		log.Printf("[ForgotPassword] BindJSON error: %v\n", err)
		if strings.Contains(err.Error(), "Email") &&
			strings.Contains(err.Error(), "email") {
			response.Error(
				ctx,
				http.StatusBadRequest,
				"Email format is invalid",
			)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	if err := c.authService.ForgotPassword(ctx.Request.Context(), user.Email); err != nil {
		log.Printf("[ForgotPasswordService] error: %v\n", err)
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, "success to sent a link to email", nil)
}

// ResetPasswordRequest godoc
//
//	@Summary		Reset password
//	@Description	Reset password using a reset token from the forgot password email.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			token	query	string					true	"Reset password token"
//	@Param			payload	body	dto.ResetPasswordBody	true	"Reset password payload"
//	@Success		200		{object}	dto.EmptyDataResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Router			/auth/password/reset [post]
func (c *AuthController) ResetPasswordRequest(ctx *gin.Context) {
	var reqBody dto.ResetPasswordBody
	var reqParam dto.ResetPasswordQuery

	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		log.Printf("[ResetPassword] JSON bind error: %v\n", err)

		if strings.Contains(err.Error(), "NewPassword") {
			if strings.Contains(err.Error(), "required") {
				response.Error(ctx, http.StatusBadRequest, "new_password is required")
				return
			}
			if strings.Contains(err.Error(), "min") {
				response.Error(ctx, http.StatusBadRequest, "Password must be at least 8 characters")
				return
			}
		}

		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := ctx.ShouldBindQuery(&reqParam); err != nil {
		log.Println("TOKEN RAW:", ctx.Query("token"))
		log.Printf("[ResetPassword] Query bind error: %v\n", err)

		if strings.Contains(err.Error(), "Token") {
			response.Error(ctx, http.StatusBadRequest, "token is required")
			return
		}

		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.authService.ResetPassword(
		ctx.Request.Context(),
		reqParam.Token,
		reqBody.NewPassword,
	); err != nil {
		log.Printf("[ResetPassword] error: %v\n", err)
		response.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(ctx, http.StatusOK, "Password reset successful", nil)
}
