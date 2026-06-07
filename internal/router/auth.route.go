package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tickitz-backend/internal/controller"
	"github.com/tickitz-backend/internal/repository"
	"github.com/tickitz-backend/internal/service"
)

func RegisterAuthRouter(router *gin.Engine, db *pgxpool.Pool) {
	authRouter := router.Group("/auth")
	authRepo := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepo)
	authController := controller.NewAuthController(authService)

<<<<<<< HEAD
	authRouter.POST("/signup", authController.Register)
	authRouter.POST("/activate", authController.Activate)
	authRouter.POST("/otp", authController.GetNewOTP)
=======
	authRouter.POST("/register", authController.Register)
>>>>>>> b9ee6f3b7daa7e17199dec072791cf7dbe5d369b
}
