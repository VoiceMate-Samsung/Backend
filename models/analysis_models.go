package models

type Move struct {
	Move string `db:"move"`
	Fen  string `db:"fen"`
}

type Game struct {
	GameID     string `db:"game_id"`
	Date       string `db:"date"`
	MoveAmount int    `db:"move_amount"`
}

type StockfishAnalysisResult struct {
	Fen      string
	BestMove string
}

const (
	BotLevelEasy   = 2
	BotLevelMedium = 5
	BotLevelhard   = 10
)
