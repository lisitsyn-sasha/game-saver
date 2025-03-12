package delete

import (
	"errors"
	"game-saver/internal/logger"
	"game-saver/internal/server/response"
	"game-saver/internal/storage"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

func NewDeleteHandler(log *slog.Logger, gameDeleter GameDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log = log.With(
			slog.String("op", "NewDeleteHandler"),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", logger.Err(err))
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			var validateErr validator.ValidationErrors
			if errors.As(err, &validateErr) {
				log.Error("invalid request", logger.Err(err))
				render.JSON(w, r, response.ValidationError(validateErr))
				return
			}

			log.Error("unexpected validation error", logger.Err(err))
			render.JSON(w, r, response.Error("invalid request"))
			return
		}

		id, err := gameDeleter.DeleteGame(r.Context(), req.Title)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("game already exists", slog.String("game", req.Title))
			render.JSON(w, r, response.Error("game title already exists"))
			return
		}
		if err != nil {
			log.Error("failed to delete game", logger.Err(err))
			render.JSON(w, r, response.Error("failed to delete game"))
			return
		}

		log.Info("game deleted", slog.Int64("id", id))

		render.JSON(w, r, response.Response{
			Status: response.OK().Status,
		})
	}
}
