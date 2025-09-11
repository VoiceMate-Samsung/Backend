package pg_sql

var (
	Move = `
	INSERT INTO public.moves (user_id, game_id, move, fen)
		VALUES ($1, $2, $3, $4)
	`
)
