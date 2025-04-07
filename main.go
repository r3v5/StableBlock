package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/r3v5/stableblock-api/database"
	"github.com/r3v5/stableblock-api/handlers"
	"github.com/r3v5/stableblock-api/middleware"
)

func main() {
    err := godotenv.Load()

    if err != nil {
		log.Fatal("❌ Error loading .env file")
	}

    database.Connect()
    
    api := gin.Default()

    api.GET("api/v1/welcome", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Welcome to StableBlock 👋",
        })
    })
    api.POST("api/v1/register", handlers.HandlePostRegister)
    api.POST("api/v1/login", handlers.HandlePostLogin)
    api.POST("api/v1/refresh", handlers.HandlePostRefresh)
    api.POST("api/v1/logout", middleware.JwtAuthMiddleware(), handlers.HandlerPostLogout)
	api.GET("api/v1/accounts", middleware.JwtAuthMiddleware(), handlers.HandleGetAccount)
    api.POST("api/v1/transactions", middleware.JwtAuthMiddleware(), handlers.HandlePostTransaction)
    api.GET("api/v1/transactions", handlers.HandleGetTransactions)
    api.GET("api/v1/transactions/:hash", handlers.HandleGetTransaction)
    api.GET("api/v1/blocks", handlers.HandleGetBlocks)
    api.GET("api/v1/blocks/:height", handlers.HandleGetBlock)
    api.GET("api/v1/sbtokens", middleware.JwtAuthMiddleware(), handlers.HandleGetSBTokens)
    api.Run()
}
