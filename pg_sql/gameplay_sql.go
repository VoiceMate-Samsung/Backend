package pg_sql

var (
	Move = `
	INSERT INTO public.moves (game_id, user_id, move, fen)
		VALUES ($1, $2, $3, $4)
	`
)
