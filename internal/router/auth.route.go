package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tickitz-backend/internal/controller"
	"github.com/tickitz-backend/internal/middleware"
	"github.com/tickitz-backend/internal/repository"
	"github.com/tickitz-backend/internal/service"
)

func RegisterAuthRouter(router *gin.Engine, db *pgxpool.Pool, authCache *repository.AuthCacheRepository) {
	authRouter := router.Group("/auth")
	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepo, authCache)
	authController := controller.NewAuthController(authService)

	authRouter.POST("/signup", authController.Register)
	authRouter.POST("/activate", authController.Activate)
	authRouter.POST("/otp", authController.GetNewOTP)
	authRouter.POST("/signin", authController.Login)
	authRouter.DELETE("/logout", middleware.VerifyToken(authCache), authController.Logout)
	authRouter.PATCH("/password", authController.ChangeUserPassword)
	authRouter.POST("/password/forgot", authController.ForgotPassword)
	authRouter.POST("/password/reset", authController.ResetPasswordRequest)
}
