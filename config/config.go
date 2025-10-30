// Package config пакет с инициализацией конфига
package config

import (
	"flag"
	"sync"

	"github.com/caarlos0/env/v11"
)

var once sync.Once

// Config структура
type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	Host          string `env:"BASE_URL"`
	DatabaseDsn   string `env:"DATABASE_DSN"`
}

// NewConfig конструктор для конфига
func NewConfig() (Config, error) {
	var conf Config

	err := env.Parse(&conf)
	if err != nil {
		return Config{}, err
	}

	if conf.Host != "" && conf.ServerAddress != "" && conf.DatabaseDsn != "" {
		return conf, nil
	}

	once.Do(func() {
		flag.StringVar(&conf.DatabaseDsn, "d", "", "database dsn") //"postgres://postgres:qwerty12345@localhost:5432/postgres"
		flag.StringVar(&conf.ServerAddress, "a", "localhost:8080", "server adress")
		flag.StringVar(&conf.Host, "b", "http://localhost:8080", "host")

		flag.Parse()
	})

	return conf, nil
}
