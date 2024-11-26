package app

import (
	"context"
	"fmt"

	"github.com/cutlery47/music-storage/internal/config"
	v1 "github.com/cutlery47/music-storage/internal/controller/http/v1"
	"github.com/cutlery47/music-storage/internal/repository"
	"github.com/cutlery47/music-storage/internal/service"
	"github.com/cutlery47/music-storage/internal/utils"
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

	infoFd, err := utils.CreateAndOpen(config.InfoPath)
	if err != nil {
		return fmt.Errorf("error when creaing info log file: %v", err)
	}

	errFd, err := utils.CreateAndOpen(config.ErrorPath)
	if err != nil {
		return fmt.Errorf("error when creating error log file: %v", err)
	}

	infoLog := logger.WithFormat(logger.WithFile(logger.New(logrus.InfoLevel), infoFd), &logrus.JSONFormatter{})
	errLog := logger.WithFormat(logger.WithFile(logger.New(logrus.ErrorLevel), errFd), &logrus.JSONFormatter{})

	repo, err := repository.NewMusicRepository(ctx, config.PostgresConfig)
	if err != nil {
		return fmt.Errorf("error when connecting to the db: %v", err)
	}

	logrus.Debug("initializing service...")
	srv := service.NewMusicService(repo)

	logrus.Debug("initializing controller...")
	echo := echo.New()
	v1.NewController(echo, srv, infoLog, errLog)

	logrus.Debug("initializing http server...")
	httpserver := httpserver.New(
		echo,
		httpserver.Addr(config.Interface, config.Port),
		httpserver.ReadTimeout(config.ReadTimeout),
		httpserver.WriteTimeout(config.WriteTimeout),
		httpserver.ShutdownTimeout(config.ShutdownTimeout),
	)

	return httpserver.Run(ctx)
}
