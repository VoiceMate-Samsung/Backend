package routes

import (
	"samsungvoicebe/config"
	"samsungvoicebe/controllers"

	"github.com/gin-gonic/gin"
)

func ChatRoutes(router *gin.RouterGroup, cfg *config.Config) {
	chatController := controllers.NewChatController(cfg)

	router.POST("/chat", chatController.Chat)
}
