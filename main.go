package main

import (
	"log"
	"subscription/api"
	database "subscription/db"

	_ "subscription/docs"

	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Subscription API
// @version 1.0
// @description Апи агрегации подписок.
func main() {
	// загружаем .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system ENV")
	}

	// инициализация БД
	database.InitDB()

	router := api.SetupRouter()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run("0.0.0.0:8081")
}
