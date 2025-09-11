package controllers

import (
	"github.com/gin-gonic/gin"
	"samsungvoicebe/config"
	"samsungvoicebe/services"
)

type AnalysisController struct {
	Config  *config.Config
	Service *services.AnalysisService
}

func NewAnalysisController(cfg *config.Config, service *services.AnalysisService) *AnalysisController {
	return &AnalysisController{
		Config:  cfg,
		Service: service,
	}
}

func (ac *AnalysisController) GetGameHistoryList(c *gin.Context) {
	userID := c.Param("user_id")
	gamesHistoryList, err := ac.Service.GetGameHistoryList(userID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": gamesHistoryList})
}
