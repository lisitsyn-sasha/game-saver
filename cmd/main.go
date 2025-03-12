package main

import (
	"game-saver/config"
	"game-saver/internal/logger"
	"game-saver/internal/server/handlers/delete"
	"game-saver/internal/server/handlers/save"
	"game-saver/internal/storage/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	// загружаем конфигурацию приложения
	cfg := config.MustLoad()

	// настраиваем логгер в зависимости от окружения
	log := logger.SetupLogger(cfg.Env)
	// записываем информационное сообщение о запуске приложения
	log.Info("starting game-saver", slog.String("env", cfg.Env))
	// записываем отладочное сообщение
	log.Debug("debug message are enabled")

	// Инициализируем хранилище, передавая путь к хранилищу из конфигурации; если есть ошибка, то выдаем ее и завершаем программу
	storage, err := postgres.NewPostgresStorage(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", logger.Err(err))
		os.Exit(1)
	}

	// Создаем новый роутер
	router := chi.NewRouter()

	router.Use(middleware.RequestID) // Добавляем уникальный id к каждому запросу
	router.Use(middleware.Logger)    // Логируем информацию о каждом запроса
	router.Use(middleware.Recoverer) // Отлавливаем паники
	router.Use(middleware.URLFormat) // Парсим и обрабатываем формат URL

	// Регистрируем маршруты и создаем обработчики
	router.Post("/game", save.NewSaveHandler(log, storage))
	router.Delete("/game", delete.NewDeleteHandler(log, storage))

	// Логируем запуск сервера
	log.Info("starting server", slog.String("address", cfg.Address))

	// Настраиваем HTTP-сервер
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	// Запускаем сервер
	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	// Логируем остановку сервера
	log.Error("server stopped")
}
