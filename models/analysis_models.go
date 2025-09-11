package models

type Move struct {
	Move string `db:"move"`
	Fen  string `db:"fen"`
}

type Game struct {
	GameID     string `db:"game_id"`
	Date       string `db:"date"`
	Fen        string `db:"fen"`
	MoveAmount int    `db:"move_amount"`
	Result     string `db:"result"`
	EndType    string `db:"end_type"`
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
