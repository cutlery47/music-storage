package v1

import (
	"github.com/cutlery47/music-storage/internal/service"
	"github.com/labstack/echo/v4"
)

type songRoutes struct {
	srv service.Service
	e   *errMapper
}

func newSongRoutes(g *echo.Group, srv service.Service, e *errMapper) {
	r := &songRoutes{
		srv: srv,
		e:   e,
	}

	g.POST("/", r.uploadSong)
	g.GET("/", r.getSongs)
	g.GET("/info", r.getInfo)
	g.GET("/text", r.getText)
	g.DELETE("/", r.deleteSong)
	g.PUT("/", r.updateSong)
}

func (r *songRoutes) getInfo(c echo.Context) error {
	return nil
}

func (r *songRoutes) getSongs(c echo.Context) error {
	return nil
}

func (r *songRoutes) getText(c echo.Context) error {
	return nil
}

func (r *songRoutes) deleteSong(c echo.Context) error {
	return nil
}

func (r *songRoutes) uploadSong(c echo.Context) error {
	return nil
}

func (r *songRoutes) updateSong(c echo.Context) error {
	return nil
}
