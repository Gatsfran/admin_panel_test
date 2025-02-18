package config

import (
	"fmt"
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
