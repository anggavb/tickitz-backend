package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/tickitz-backend/internal/config"
	"github.com/tickitz-backend/internal/router"
)

// @title						Backend Tickitz API
// @version						1.0
// @description					API documentation for Tickitz backend application

// @license.name			MIT

// @host					localhost:8081
// @BasePath				/

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description					Bearer token used for authorization. Example: Bearer <token>
func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading env. \ncause: %s", err.Error())
	}

	app := gin.Default()

	db, err := config.ConnectPsql()
	if err != nil {
		log.Printf("DB connection error. \ncause: %s", err.Error())
	}
	defer db.Close()
	log.Println("DB Connected")

	rdb, err := config.ConnectRDB()
	if err != nil {
		log.Printf("Redis connection error. \ncause: %s", err.Error())
	}
	defer rdb.Close()
	log.Println("Redis Connected")

	router.InitRouter(app, db, rdb)
	app.Run(fmt.Sprintf("%s:%s", os.Getenv("APP_HOST"), os.Getenv("APP_PORT")))

}
