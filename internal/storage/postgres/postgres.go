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
            title TEXT NOT NULL UNIQUE
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
	INSERT INTO game(title) VALUES($1)
	ON CONFLICT (title) DO UPDATE SET title = EXCLUDED.title
	RETURNING id
`

// SaveGame сохраняет игру в БД и возвращает ее id
// Если gameToSave пустое, то возвращает ошибку
func (s *Storage) SaveGame(ctx context.Context, gameTitleToSave string) (int64, error) {
	// Проверяем, что gameToSave не пустое
	if gameTitleToSave == "" {
		return 0, fmt.Errorf("game title cannot be empty")
	}

	var id int64

	// Выполняет SQL-запрос; если произошла ошибка, то возвращаем ее с контекстом
	err := s.db.QueryRow(ctx, saveGameQuery, gameTitleToSave).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to save game: %w", err)
	}

	// возвращаем id сохраненной игры
	return id, nil
}
