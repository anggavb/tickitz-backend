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
