package save

import "context"

type GameSaver interface {
	SaveGame(ctx context.Context, gameTitleToSave string) (int64, error)
}

type Request struct {
	Title string `json:"title" validate:"required"`
}
