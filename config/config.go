package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

// MustLoad нужна для загрузки конфигурации в начале работы приложения и отслеживает ранние ошибки
func MustLoad() *Config {
	// получаем значение переменной окружения CONFIG_PATH, которая содержит путь к файлу local.yaml
	configPath := os.Getenv("CONFIG_PATH")

	// проверяем, установлена ли переменная окружения CONFIG_PATH; если нет, то завершаем приложение с ошибкой
	if configPath == "" {
		log.Fatalf("CONFIG_PATH is not set")
	}

	// проверяем, существует ли файл по пути configPath; если нет, то завершаем приложение с ошибкой
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	// создаем переменную cfg типа Config; это структура, которая будет хранить конфигурацию
	var cfg Config

	// используем библиотеку cleanenv для чтения конфигурации и заполнения структуры cfg
	// если конфигурация имеет неправильный формат или содержит ошибки, программа завершится с ошибкой
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	// возвращаем указатель на структуру cfg, которая содержит загруженную конфигурацию
	return &cfg
}
