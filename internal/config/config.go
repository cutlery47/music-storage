package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	HttpConfig
	PostgresConfig
	LoggerConfig
}

type HttpConfig struct {
	Port            string        `env:"HTTP_PORT"`
	Interface       string        `env:"HTTP_INTERFACE"`
	ReadTimeout     time.Duration `env:"HTTP_READ_TIMEOUT"`
	WriteTimeout    time.Duration `env:"HTTP_WRITE_TIMEOUT"`
	ShutdownTimeout time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT"`
}

type PostgresConfig struct {
	PostgresUser     string `env:"POSTGRES_USER"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"`
	PostgresPort     string `env:"POSTGRES_PORT"`
	PostgresSSL      string `env:"POSTGRES_SSL"`
}

type LoggerConfig struct {
	InfoPath  string `env:"INFO_LOGS_PATH"`
	ErrorPath string `env:"ERROR_LOGS_PATH"`
}

func New() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, fmt.Errorf("godotenv.Load: %v", err)
	}

	var (
		pgConf   PostgresConfig
		logConf  LoggerConfig
		httpConf HttpConfig
	)

	if err := cleanenv.ReadEnv(&pgConf); err != nil {
		return nil, fmt.Errorf("couldn't read postgres config: %v", err)
	}

	if err := cleanenv.ReadEnv(&httpConf); err != nil {
		return nil, fmt.Errorf("couldn't read http config: %v", err)
	}

	if err := cleanenv.ReadEnv(&logConf); err != nil {
		return nil, fmt.Errorf("coundn't read logger config: %v", err)
	}

	return &Config{
		PostgresConfig: pgConf,
		HttpConfig:     httpConf,
		LoggerConfig:   logConf,
	}, nil
}
