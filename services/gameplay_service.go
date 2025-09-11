package services

import (
	"fmt"

	"samsungvoicebe/models"
	"samsungvoicebe/repo"
)

type GameplayService struct {
	gameplayRepo    *repo.GameplayRepo
	analysisService *AnalysisService
}

func NewGameplayService(gameplayRepo *repo.GameplayRepo, analysisService *AnalysisService) *GameplayService {
	return &GameplayService{
		gameplayRepo:    gameplayRepo,
		analysisService: analysisService,
	}
}

func (s *GameplayService) PlayerMove(userID, gameID, fen, move, botLevel string) (models.BotMove, error) {
	err := s.gameplayRepo.GameMove(userID, gameID, fen, move)
	if err != nil {
		err = fmt.Errorf("GameplayService-PlayerMove-GameMove: %w", err)
		fmt.Printf("userID: %s, gameID: %s, fen: %s, move: %s", userID, gameID, fen, move)
		return models.BotMove{}, err
	}

	analysisResult, err := s.analysisService.StockfishAnalyze(fen, botLevel)
	if err != nil {
		err = fmt.Errorf("GameplayService-PlayerMove-GameMove: %w", err)
		return models.BotMove{}, err
	}

	var botMove models.BotMove

	botMove = models.BotMove{
		Fen:  analysisResult.Fen,
		Move: analysisResult.BestMove,
	}

	return botMove, nil
}
