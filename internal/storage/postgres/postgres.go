package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

// NewPostgresStorage создает подключение к БД и инициализирует таблицы
func NewPostgresStorage(connString string) (*Storage, error) {
	// создаем пул соединения с БД; если создать не удалось, возвращаем ошибку
	db, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// создаем таблицы в БД; если произошла ошибка, то закрываем соединение и возвращаем ошибку
	err = createTables(db)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	// возвращаем структуру Storage с подключением к БД
	return &Storage{db: db}, nil
}

// Close закрывает пул соединений с БД
func (s *Storage) Close() {
	s.db.Close()
}

// объявляем переменные для хранения SQL-запросов
var (
	createTableURLQuery = `
        CREATE TABLE IF NOT EXISTS game(
            id SERIAL PRIMARY KEY,
            title TEXT NOT NULL UNIQUE,
            score INTEGER
        );
    `
	createIdURLQuery = `CREATE INDEX IF NOT EXISTS idx_game_title ON game(title);`
)

// createTables создаем таблицу и индекс в БД
// createTables принимает пул соединений с БД и возвращает ошибку
func createTables(db *pgxpool.Pool) error {
	// выполняем sql-запрос для создания таблицы `game`
	_, err := db.Exec(context.Background(), createTableURLQuery)
	if err != nil {
		return fmt.Errorf("failed to create game table: %w", err)
	}

	// выполняем sql-запрос для создания индекса `idx_game_title` на столбце `title`
	_, err = db.Exec(context.Background(), createIdURLQuery)
	if err != nil {
		return fmt.Errorf("failed to create index on game: %w", err)
	}

	// если оба запроса выполнены успешно, возвращаем nil
	return nil
}

// SQL-запрос для сохранения игры
// Если игра с таким названием существует, обновляется запись и возвращается id
var saveGameQuery = `
	INSERT INTO game(title, score) VALUES($1, $2)
	ON CONFLICT (title) DO UPDATE SET score = EXCLUDED.score
	RETURNING id
`

// SaveGame сохраняет игру в БД и возвращает ее id
// Если gameToSave пустое, то возвращает ошибку
func (s *Storage) SaveGame(ctx context.Context, gameTitleToSave string, gameScoreToSave int64) (int64, error) {
	// Проверяем, что gameToSave не пустое
	if gameTitleToSave == "" {
		return 0, fmt.Errorf("game title cannot be empty")
	}
	// Проверяем, что оценка больше 0 и меньше 10
	if gameScoreToSave < 0 || gameScoreToSave > 10 {
		return 0, fmt.Errorf("score must be between 0 and 10")
	}

	var id int64

	// Выполняет SQL-запрос; если произошла ошибка, то возвращаем ее с контекстом
	err := s.db.QueryRow(ctx, saveGameQuery, gameTitleToSave, gameScoreToSave).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to save game: %w", err)
	}

	// возвращаем id сохраненной игры
	return id, nil
}

// SQL-запрос с удалением игры
const deleteGameQuery = `DELETE FROM game WHERE title = $1 RETURNING id`

// DeleteGame удаляет игру из БД; функция написана аналогично SaveGame
func (s *Storage) DeleteGame(ctx context.Context, gameTitleToDelete string) (int64, error) {
	if gameTitleToDelete == "" {
		return 0, fmt.Errorf("game title cannot be empty")
	}

	var id int64

	err := s.db.QueryRow(ctx, deleteGameQuery, gameTitleToDelete).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to delete game: %w", err)
	}

	return id, nil
}

// SQL-запрос на получение оценки игры
const getGameScoreQuery = `SELECT score FROM game WHERE title = $1`

// GetGameScore запрашивает оценку игры из БД; функция написана аналогично SaveGame
func (s *Storage) GetGameScore(ctx context.Context, gameTitle string) (int64, error) {
	if gameTitle == "" {
		return 0, fmt.Errorf("game title cannot be empty")
	}

	var score int64

	err := s.db.QueryRow(ctx, getGameScoreQuery, gameTitle).Scan(&score)
	if err != nil {
		return 0, fmt.Errorf("failed to get score: %w", err)
	}

	return score, nil
}
