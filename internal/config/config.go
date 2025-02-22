package config

import (
	"fmt"
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Postgres struct {
	DSN      string `env:"POSTGRES_DSN,required"`
	Host     string `env:"HOST" envDefault:"localhost"`
	Port     string `env:"PORT" envDefault:"15432"`
	Username string `env:"USERNAME" envDefault:"postgres"`
	Password string `env:"PASSWORD, required"`
	Database string `env:"NAME" envDefault:"admin_db"`
}
type Server struct {
	Port string `env:"PORT" envDefault:"8080"`
}

type Telegram struct {
	Token  string `env:"TELEGRAM_TOKEN,required"`
	ChatID int64  `env:"TELEGRAM_CHAT_ID,required"`
}
type Config struct {
	Postgres      Postgres      `envPrefix:"POSTGRES_"`
	Server        Server        `envPrefix:"SERVER_"`
	JWTSecret     string        `env:"JWT_SECRET,required"`
	JWTExpiration time.Duration `env:"JWT_EXPIRATION,required"`
	Telegram      Telegram
}

func New() *Config {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("Ошибка при парсинге переменных окружения: %v", err)
	}
	fmt.Println(cfg)
	return &cfg
}

func (p *Postgres) GetPostgresConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		p.Username,
		p.Password,
		p.Host,
		p.Port,
		p.Database,
	)
}
