package services

import (
	"fmt"

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
	err := s.gameplayRepo.GameMove(userID, gameID, fen, move)
	if err != nil {
		fmt.Errorf("GameplayService-PlayerMove-GameMove: %w", err)
		return models.BotMove{}, err
	}

	var botMove models.BotMove
	return botMove, nil
}
