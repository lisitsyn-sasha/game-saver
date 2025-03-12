package delete

import "context"

type GameDeleter interface {
	DeleteGame(ctx context.Context, gameTitleToDelete string) (int64, error)
}

type Request struct {
	Title string `json:"title" validate:"required"`
}
