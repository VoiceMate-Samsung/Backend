package pg_sql

var (
	GetGameHistoryList = `
	SELECT
		game_id,
		created_at AS date,
		fen,
		move_amount,
		result,
		end_type
	FROM public.games
	WHERE user_id = $1 AND game_id = $2
	ORDER BY created_at DESC
	`

	GetMoveByOrder = `
	SELECT move, fen FROM public.moves
		WHERE game_id = $1 AND user_id = $2 AND move_order = $3

	`
)
