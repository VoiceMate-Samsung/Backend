package services

import (
	"fmt"

	"github.com/notnil/chess"
	"github.com/notnil/chess/uci"
	"samsungvoicebe/models"
	"samsungvoicebe/repo"
)

type AnalysisService struct {
	analysisRepo *repo.AnalysisRepo
}

func NewAnalysisService(analysisRepo *repo.AnalysisRepo) *AnalysisService {
	return &AnalysisService{analysisRepo: analysisRepo}
}

func (a *AnalysisService) StockfishAnalyze(fen string, botLevel string) (models.StockfishAnalysisResult, error) {
	var analysisResult models.StockfishAnalysisResult

	engine, err := uci.New("stockfish")
	if err != nil {
		err = fmt.Errorf("AnalysisService-StockfishAnalyze-uci.New: %w", err)
		return models.StockfishAnalysisResult{}, err
	}
	defer engine.Close()

	err = engine.Run(uci.CmdUCI, uci.CmdIsReady)
	if err != nil {
		err = fmt.Errorf("AnalysisService-StockfishAnalyze-engine.Run-uci.CmdUCI-uci.CmdIsReady: %w", err)
		return models.StockfishAnalysisResult{}, err
	}

	position, err := chess.FEN(fen)
	if err != nil {
		err = fmt.Errorf("AnalysisService-StockfishAnalyze-chess.FEN: %w", err)
		return models.StockfishAnalysisResult{}, err
	}

	game := chess.NewGame(position)

	err = engine.Run(uci.CmdPosition{Position: game.Position()})
	if err != nil {
		err = fmt.Errorf("AnalysisService-StockfishAnalyze-engine.Run-uci.CmdPosition: %w", err)
		return models.StockfishAnalysisResult{}, err
	}

	var depth int
	switch botLevel {
	case "easy":
		depth = models.BotLevelEasy
	case "medium":
		depth = models.BotLevelMedium
	case "hard":
		depth = models.BotLevelhard
	default:
		depth = models.BotLevelMedium
	}

	searchBestMove := uci.CmdGo{Depth: depth}

	err = engine.Run(searchBestMove)
	if err != nil {
		err = fmt.Errorf("AnalysisService-StockfishAnalyze-engine.Run-searchBestMove: %w", err)
		return models.StockfishAnalysisResult{}, err
	}

	bestMove := engine.SearchResults().BestMove

	err = game.Move(bestMove)
	if err != nil {
		err = fmt.Errorf("AnalysisService-StockfishAnalyze-game.Move: %w", err)
		return models.StockfishAnalysisResult{}, err
	}

	analysisResult = models.StockfishAnalysisResult{
		BestMove: bestMove.String(),
		Fen:      game.FEN(),
	}

	return analysisResult, nil
}

func (a *AnalysisService) GetGameHistoryList(userID string) ([]models.Game, error) {
	games, err := a.analysisRepo.GetGameHistoryList(userID)
	if err != nil {
		err = fmt.Errorf("AnalysisService-GetGameHistoryList-GetGameHistoryList: %w", err)
		return []models.Game{}, err
	}
	
	return games, nil
}
