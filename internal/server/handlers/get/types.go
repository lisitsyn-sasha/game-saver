package get

import "context"

type GameGetter interface {
	GetGameScore(ctx context.Context, gameTitle string) (int64, error)
}

type Request struct {
	Title string `json:"title" validate:"required"`
}
