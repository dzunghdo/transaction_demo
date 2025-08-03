package config

import "fmt"

// Config represents the application configuration
type Config struct {
	AppName  string   `mapstructure:"app_name"`
	Env      string   `mapstructure:"env"`
	Server   Server   `mapstructure:"server"`
	Postgres Postgres `mapstructure:"postgres"`
}

type Server struct {
	Port uint `mapstructure:"port"`
}

type Postgres struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DB           string `mapstructure:"db"`
	Port         string `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

func (p *Postgres) Conn() string {
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		p.User, p.Password, p.Host, p.Port, p.DB)
}
