package v1

import (
	"github.com/cutlery47/music-storage/internal/service"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func NewController(e *echo.Echo, srv service.Service, infoLog, errLog *logrus.Logger) {
	// healthcheck endpoing
	e.GET("/ping", func(c echo.Context) error { return c.NoContent(200) })
	// swagger endpoint
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	v1 := e.Group("/api/v1/songs", requestLoggerMiddleware(infoLog))
	{
		newSongRoutes(v1, srv, newErrMapper(errLog))
	}

}
