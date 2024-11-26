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

type Mode struct {
	Mode string `env:"APP_MODE"`
}

type HttpConfig struct {
	Port            string        `env:"HTTP_PORT"`
	Interface       string        `env:"HTTP_INTERFACE"`
	ReadTimeout     time.Duration `env:"HTTP_READ_TIMEOUT"`
	WriteTimeout    time.Duration `env:"HTTP_WRITE_TIMEOUT"`
	ShutdownTimeout time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT"`
}

type PostgresConfig struct {
	PostgresUser       string        `env:"POSTGRES_USER"`
	PostgresPassword   string        `env:"POSTGRES_PASSWORD"`
	PostgresHost       string        `env:"POSTGRES_HOST"`
	PostgresPort       string        `env:"POSTGRES_PORT"`
	PostgresDB         string        `env:"POSTGRES_DB"`
	PostgresSSL        string        `env:"POSTGRES_SSL"`
	PostgresMigrations string        `env:"POSTGRES_MIGRATIONS_PATH"`
	PostgresTimeout    time.Duration `env:"POSTGRES_CONN_TIMEOUT"`
}

type LoggerConfig struct {
	InfoPath  string `env:"INFO_LOGS_PATH"`
	ErrorPath string `env:"ERROR_LOGS_PATH"`
}

func New() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, fmt.Errorf("godotenv.Load: %v", err)
	}

	mode := Mode{}
	if err := cleanenv.ReadEnv(&mode); err != nil {
		return nil, fmt.Errorf("couldn't read mode config: %v", err)
	}

	var (
		pgConf   PostgresConfig
		logConf  LoggerConfig
		httpConf HttpConfig
	)

	switch mode.Mode {
	case "DEV":
		setDevConfig(&pgConf, &logConf, &httpConf)
	case "PROD":
		if err := setProdConfig(&pgConf, &logConf, &httpConf); err != nil {
			return nil, fmt.Errorf("setProdConfig: %v", err)
		}
	default:
		return nil, fmt.Errorf("only DEV and PROD modes are allowed...")
	}

	return &Config{
		PostgresConfig: pgConf,
		HttpConfig:     httpConf,
		LoggerConfig:   logConf,
	}, nil
}

func setProdConfig(pgConf *PostgresConfig, logConf *LoggerConfig, httpConf *HttpConfig) error {
	if err := cleanenv.ReadEnv(pgConf); err != nil {
		return fmt.Errorf("couldn't read postgres config: %v", err)
	}

	if err := cleanenv.ReadEnv(httpConf); err != nil {
		return fmt.Errorf("couldn't read http config: %v", err)
	}

	if err := cleanenv.ReadEnv(logConf); err != nil {
		return fmt.Errorf("coundn't read logger config: %v", err)
	}

	return nil
}

func setDevConfig(pgConf *PostgresConfig, logConf *LoggerConfig, httpConf *HttpConfig) {
	pgConf.PostgresDB = "music"
	pgConf.PostgresHost = "localhost"
	pgConf.PostgresPort = "5432"
	pgConf.PostgresUser = "postgres"
	pgConf.PostgresPassword = "12345"
	pgConf.PostgresSSL = "disable"
	pgConf.PostgresMigrations = "migrations/v2"
	pgConf.PostgresTimeout = 3 * time.Second

	logConf.ErrorPath = "logs/err.log"
	logConf.InfoPath = "logs/info.log"

	httpConf.Port = "8080"
	httpConf.Interface = "0.0.0.0"
	httpConf.ReadTimeout = 3 * time.Second
	httpConf.WriteTimeout = 3 * time.Second
	httpConf.ShutdownTimeout = 3 * time.Second
}
