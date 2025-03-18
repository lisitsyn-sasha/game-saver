package get

import (
	"game-saver/internal/logger"
	"game-saver/internal/server/response"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

// Более подробные комментарии можно найти в NewSaveHandler
func NewGetScoreHandler(log *slog.Logger, gameGetter GameGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(
			slog.String("op", "NewGetScoreHandler"),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// Получаем title из query-параметра URL
		// Если параметр отсутствует, логируем ошибку и возвращаем JSON-ответ с ошибкой
		title := r.URL.Query().Get("title")
		if title == "" {
			log.Error("missing title parameter")
			render.JSON(w, r, response.Error("title is required"))
			return
		}

		log.Info("request received", slog.String("title", title))

		score, err := gameGetter.GetGameScore(r.Context(), title)
		if err != nil {
			log.Error("failed to get game score", logger.Err(err))
			render.JSON(w, r, response.Error("failed to get game score"))
			return
		}

		log.Info("game score retrieved", slog.Int64("score", score))

		render.JSON(w, r, response.Response{
			Status: response.OK().Status,
			Data:   map[string]int64{"score": score},
		})
	}
}
