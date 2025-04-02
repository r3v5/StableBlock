package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/r3v5/stableblock-api/database"
)

func main() {
    err := godotenv.Load()

    if err != nil {
		log.Fatal("❌ Error loading .env file")
	}

    database.Connect()
	database.Migrate()
    
    r := gin.Default()

    // Basic route
    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "Welcome to StableBlock 👋",
        })
    })

    r.Run() // default: localhost:8080
}
