package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
<<<<<<< HEAD
)

func InitRouter(router *gin.Engine, db *pgxpool.Pool) {
=======
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/tickitz-backend/docs"
	"github.com/tickitz-backend/internal/middleware"
)

func InitRouter(router *gin.Engine, db *pgxpool.Pool) {
	router.Use(middleware.CORSMiddleware)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

>>>>>>> b9ee6f3b7daa7e17199dec072791cf7dbe5d369b
	RegisterAuthRouter(router, db)
}
