package services

import (
	"samsungvoicebe/models"
	"samsungvoicebe/repo"
)

type GameplayService struct {
	gameplayRepo *repo.GameplayRepo
}

func NewGameplayService(gameplayRepo *repo.GameplayRepo) *GameplayService {
	return &GameplayService{gameplayRepo: gameplayRepo}
}

func (s *GameplayService) PlayerMove(userID, gameID, fen, move string) (models.BotMove, error) {
	// TODO: continue from here to make the bot move and stuff
	var botMove models.BotMove
	return botMove, nil
}
