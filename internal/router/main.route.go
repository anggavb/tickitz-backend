package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/tickitz-backend/docs"
	"github.com/tickitz-backend/internal/middleware"
	"github.com/tickitz-backend/internal/repository"
)

func InitRouter(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	router.Use(middleware.CORSMiddleware)
	router.Static("/img", "./public/img")

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	authCache := repository.NewAuthCacheRepository(rdb)

	RegisterAuthRouter(router, db, authCache)
	RegisterMovieRouter(router, db, authCache)
	RegisterProfileRouter(router, db, authCache)
	HomeMovieRouter(router, db)
	RegisterSeatRouter(router, db)
	RegisterOrderRouter(router, db)
}
