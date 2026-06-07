package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
<<<<<<< HEAD
	_ "github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
=======
	"github.com/joho/godotenv"
>>>>>>> b9ee6f3b7daa7e17199dec072791cf7dbe5d369b
	"github.com/tickitz-backend/internal/config"
	"github.com/tickitz-backend/internal/router"
)

// @title						Backend Tickitz API
// @version						1.0
// @description					API documentation for Tickitz backend application

// @license.name				MIT

// @host						localhost:8080
// @BasePath					/

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description					Bearer token used for authorization. Example: Bearer <token>
func main() {
	fmt.Println("masuk")
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading env. \ncause: %s", err.Error())
	}

	app := gin.Default()

	db, err := config.ConnectPsql()
	if err != nil {
		log.Fatalf("DB connection error. \ncause: %s", err.Error())
	}
	defer db.Close()
	log.Println("DB Connected")
	router.InitRouter(app, db)
	app.Run(fmt.Sprintf("%s:%s", os.Getenv("APP_HOST"), os.Getenv("APP_PORT")))

}
