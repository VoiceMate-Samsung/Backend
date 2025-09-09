package models

type BotMove struct {
	Move       string
	Fen        string
	HasEnd     bool
	EndMessage string
}
