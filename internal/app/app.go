package app

import (
	"context"
	"fmt"

	"github.com/cutlery47/music-storage/internal/config"
	v1 "github.com/cutlery47/music-storage/internal/controller/http/v1"
	"github.com/cutlery47/music-storage/internal/repository"
	"github.com/cutlery47/music-storage/internal/service"
	"github.com/cutlery47/music-storage/pkg/httpserver"
	"github.com/cutlery47/music-storage/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// @title           Online Music Storage Service
// @version         0.0.1
// @description     This a service for storing music

// @contact.name   Ivanchenko Arkhip
// @contact.email  kitchen_cutlery@mail.ru

// @BasePath  /

func Run() error {
	ctx := context.Background()

	config, err := config.New()
	if err != nil {
		return fmt.Errorf("error when parsing config: %v", err)
	}

	debugLog, err := logger.NewDefaultCli(logrus.DebugLevel)
	if err != nil {
		return fmt.Errorf("error when creating debug logger: %v", err)
	}

	infoLog, err := logger.NewJsonFile(config.InfoPath, logrus.InfoLevel)
	if err != nil {
		return fmt.Errorf("error when creating info logger: %v", err)
	}

	errLog, err := logger.NewJsonFile(config.ErrorPath, logrus.ErrorLevel)
	if err != nil {
		return fmt.Errorf("error when creating debug loffer: %v", err)
	}

	url := fmt.Sprintf(
		"postgresql://%v:%v@localhost:%v/music?sslmode=%v",
		config.PostgresUser,
		config.PostgresPassword,
		config.PostgresPort,
		config.PostgresSSL,
	)

	debugLog.Debug("initializing repo...")
	repo, err := repository.NewMusicRepository(url)
	if err != nil {
		return fmt.Errorf("error when connecting to the db: %v", err)
	}

	debugLog.Debug("initializing service...")
	srv := service.NewMusicService(repo)

	debugLog.Debug("initializing controller...")
	echo := echo.New()
	v1.NewController(echo, srv, infoLog, errLog)

	debugLog.Debug("initializing http server...")
	httpserver := httpserver.New(
		echo,
		httpserver.Addr(config.Interface, config.Port),
		httpserver.ReadTimeout(config.ReadTimeout),
		httpserver.WriteTimeout(config.WriteTimeout),
		httpserver.ShutdownTimeout(config.ShutdownTimeout),
	)

	return httpserver.Run(ctx, debugLog)
}
