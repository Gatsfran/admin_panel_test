package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Postgres struct {
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	Port     string `env:"DB_PORT" envDefault:"15423"`
	Username string `env:"DB_USERNAME" envDefault:"postgres"`
	Password string `env:"DB_PASSWORD" envDefault:"docker"`
	Database string `env:"DB_NAME" envDefault:"admin_db"`
}
type Server struct {
	Port string `env:"SERVER_PORT" envDefault:"8080"`
}
type Config struct {
	Postgres Postgres
	Server   Server
}

func New() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден, используются переменные окружения по умолчанию")
	}
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Ошибка при парсинге переменных окружения: %v", err)
	}
	return cfg
}

func (c *Config) GetPostgresConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		c.Postgres.Username,
		c.Postgres.Password,
		c.Postgres.Host,
		c.Postgres.Port,
		c.Postgres.Database,
	)
}
