package save

import "context"

type GameSaver interface {
	SaveGame(ctx context.Context, gameTitleToSave string, gameScoreToSave int64) (int64, error)
}

type Request struct {
	Title string `json:"title" validate:"required"`
	Score int64  `json:"score" validate:"required"`
}
