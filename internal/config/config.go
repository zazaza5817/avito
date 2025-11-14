// Package config обрабатывает загрузку конфигурации приложения из переменных окружения.
package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config содержит конфигурацию приложения
type Config struct {
	DB     DatabaseConfig
	Server ServerConfig
	Env    string
}

// DatabaseConfig содержит параметры подключения к БД
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// ServerConfig содержит параметры HTTP сервера
type ServerConfig struct {
	Port string
	Host string
}

// Load загружает конфигурацию из переменных окружения
func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		DB: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "pr_reviewer_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
		},
		Env: getEnv("ENV", "development"),
	}

	return cfg, nil
}

// GetDSN возвращает строку подключения к PostgreSQL
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DB.Host,
		c.DB.Port,
		c.DB.User,
		c.DB.Password,
		c.DB.DBName,
		c.DB.SSLMode,
	)
}

// GetServerAddr возвращает адрес сервера
func (c *Config) GetServerAddr() string {
	return c.Server.Host + ":" + c.Server.Port
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
