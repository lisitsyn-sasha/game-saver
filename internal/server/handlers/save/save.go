package save

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

func NewSaveHandler(log *slog.Logger, gameSaver GameSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Добавляем дополнительные поля в логгер для текущего запроса
		log = log.With(
			slog.String("op", "NewSaveHandler"),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// Создаем переменную для хранения данных из JSON-запроса
		var req Request

		// Докедируем тело запроса в структуру req
		// Если декодирование не удалось, то логируем ошибку, возвращаем JSON-ответ с ошибкой и статусов 404
		// Прекращаем выполнение функции
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", logger.Err(err))
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		// Проверяем, соответствует ли структура req заданным правилам валидации
		// Если валидация не прошла, то логируем ошибку, возвращаем JSON-ответ с ошибкой валидации и возвращаем общее сообщение об ошибке
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

		// Сохраняем игру в хранилище
		// Если игра уже существует, то логируется информация и возвращает JSON-ответ с ошибкой
		// Если произошла дугая ошибка, то она возвращается
		id, err := gameSaver.SaveGame(r.Context(), req.Title)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("game already exists", slog.String("game", req.Title))
			render.JSON(w, r, response.Error("game title already exists"))
			return
		}
		if err != nil {
			log.Error("failed to add game", logger.Err(err))
			render.JSON(w, r, response.Error("failed to add game"))
			return
		}

		// Логируем успешный запрос
		log.Info("game added", slog.Int64("id", id))

		// Возвращаем успешный ответ
		render.JSON(w, r, response.Response{
			Status: response.OK().Status,
		})
	}
}
