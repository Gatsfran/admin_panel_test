package config

import (
	"fmt"
	"log"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Postgres struct {
	DSN      string `envconfig:"POSTGRES_DSN,required"`
	Host     string `envconfig:"HOST" envDefault:"localhost"`
	Port     string `envconfig:"PORT" envDefault:"15432"`
	Username string `envconfig:"USERNAME" envDefault:"postgres"`
	Password string `envconfig:"PASSWORD" required:"true"`
	Database string `envconfig:"NAME" envDefault:"admin_db"`
}
type Server struct {
	Port         string `envconfig:"PORT" envDefault:"8080"`
	IsProduction bool   `envconfig:"IS_PRODUCTION" envDefault:"true"`
	CORS         CORS   `envconfig:"CORS"`
}

type CORS struct {
	Allow_origins []string `envconfig:"CORS_ALLOW_ORIGINS" required:"true"`
	Allow_methods []string `envconfig:"CORS_ALLOW_METHODS" required:"true"`
	Allow_headers []string `envconfig:"CORS_ALLOW_HEADERS" required:"true"`
}

type Telegram struct {
	Token  string `envconfig:"TOKEN" required:"true"`
	ChatID int64  `envconfig:"CHAT_ID" required:"true"`
}
type Config struct {
	Postgres      *Postgres     `envconfig:"POSTGRES"`
	Server        *Server       `envconfig:"SERVER"`
	JWTSecret     string        `envconfig:"JWT_SECRET" required:"true"`
	JWTExpiration time.Duration `envconfig:"JWT_EXPIRATION" required:"true"`
	Telegram      *Telegram     `envconfig:"TELEGRAM"`
}

func New() *Config {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("Ошибка при парсинге переменных окружения: %v", err)
	}
	fmt.Printf("Telegram Token: %s\n", cfg.Telegram.Token)
	fmt.Printf("Telegram ChatID: %d\n", cfg.Telegram.ChatID)
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
