package repo

import (
	"database/sql"

	"samsungvoicebe/pg_sql"
)

type GameplayRepo struct {
	db *sql.DB
}

func NewGameplayRepo(db *sql.DB) *GameplayRepo {
	return &GameplayRepo{db: db}
}

func (r *GameplayRepo) GameMove(userID, gameID, fen, move string) error {
	_, err := r.db.Exec(pg_sql.Move, userID, gameID, fen, move)
	if err != nil {
		return err
	}
	return nil
}
