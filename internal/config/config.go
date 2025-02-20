package config

import (
	"fmt"
	"log"
	

	"github.com/caarlos0/env"
	
)

type Postgres struct {
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	Port     string `env:"DB_PORT" envDefault:"15423"`
	Username string `env:"DB_USERNAME" envDefault:"postgres"`
	Password string `env:"DB_PASSWORD, required"`
	Database string `env:"DB_NAME" envDefault:"admin_db"`
}
type Server struct {
	Port string `env:"SERVER_PORT" envDefault:"8080"`
}
type Config struct {
	Postgres Postgres
	Server   Server
	JWTSecret string `env:"JWT_SECRET,required"`

}

func New() Config {

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Ошибка при парсинге переменных окружения: %v", err)
	}
	return cfg
}

func (c *Postgres) GetPostgresConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		c.Username,
		c.Password,
		c.Host,
		c.Port,
		c.Database,
	)
}

