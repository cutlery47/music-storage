package v1

import (
	"github.com/cutlery47/music-storage/internal/repository"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

var (
	ErrBadQuery           = echo.NewHTTPError(400, "all required query parameters should be provided...")
	ErrBadQueryTime       = echo.NewHTTPError(400, "couldn't parse provided time...")
	ErrBadQueryPagination = echo.NewHTTPError(400, "couldn't parse pagination params...")
)

var errMap = map[error]*echo.HTTPError{
	repository.ErrNotFound:      echo.ErrNotFound,
	repository.ErrAlreadyExists: echo.ErrBadRequest,
}

type errMapper struct {
	errLog *logrus.Logger
}

func newErrMapper(errLog *logrus.Logger) *errMapper {
	return &errMapper{
		errLog: errLog,
	}
}

func (e errMapper) Map(err error) *echo.HTTPError {
	if httpErr, ok := errMap[err]; ok {
		httpErr.Message = err.Error()
		return httpErr
	}

	e.errLog.Error(err.Error())
	return echo.ErrInternalServerError
}
