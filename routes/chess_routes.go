package routes

import (
	"samsungvoicebe/config"
	"samsungvoicebe/controllers"

	"github.com/gin-gonic/gin"
)

func ChessRoutes(router *gin.RouterGroup, cfg *config.Config) {
	chessController := controllers.NewChessController(cfg)

	router.POST("/ai", chessController.PlayChess)
}
