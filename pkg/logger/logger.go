package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func NewJsonFile(filepath string, level logrus.Level) (*logrus.Logger, error) {
	logger := logrus.New()

	fd, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	logger.SetOutput(fd)
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(level)

	return logger, nil
}

func NewDefaultCli(level logrus.Level) (*logrus.Logger, error) {
	logger := logrus.New()

	logger.SetOutput(os.Stdout)
	logger.SetLevel(level)

	return logger, nil
}
