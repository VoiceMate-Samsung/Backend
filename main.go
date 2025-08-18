package main

import (
	"fmt"
	"log"
	"samsungvoicebe/config"
	"samsungvoicebe/middleware"
	"samsungvoicebe/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	if cfg.IsValid() {
		log.Println("✅ Environment loaded successfully")
	} else {
		log.Fatal("❌ GEMINI_API_KEY not configured")
	}

	gin.SetMode(cfg.GinMode)

	r := gin.New()

	r.Use(gin.Logger())
	r.Use(middleware.CORSMiddleware())
	r.Use(gin.Recovery())

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Samsung Voice Backend API",
			"version": "1.0.0",
			"status":  "running",
		})
	})

	chatApi := r.Group("/api/chat")
	routes.ChatRoutes(chatApi, cfg)

	chessApi := r.Group("/api/chess")
	routes.ChessRoutes(chessApi, cfg)

	fmt.Printf("Base URL: http://localhost:%s/\n", cfg.Port)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("❌ Failed to start server:", err)
	}
}
