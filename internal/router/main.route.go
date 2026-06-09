package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// _ "github.com/tickitz-backend/docs"
	"github.com/tickitz-backend/internal/middleware"
)

func InitRouter(router *gin.Engine, db *pgxpool.Pool) {
	router.Use(middleware.CORSMiddleware)
	router.Static("/img", "./public/img")

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	RegisterAuthRouter(router, db)
	RegisterMovieRouter(router, db)
	HomeMovieRouter(router, db)
}
