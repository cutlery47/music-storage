package v1

import (
	"strconv"
	"time"

	"github.com/cutlery47/music-storage/internal/models"
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

	g.POST("", r.uploadSong)
	g.GET("", r.getSongs)
	g.GET("/info", r.getInfo)
	g.GET("/text", r.getText)
	g.DELETE("", r.deleteSong)
	g.PUT("", r.updateSong)
}

func (r *songRoutes) getInfo(c echo.Context) error {
	params := c.QueryParams()

	if !params.Has("group") || !params.Has("song") {
		return ErrBadQuery
	}

	song := models.Song{
		GroupName: params.Get("group"),
		SongName:  params.Get("song"),
	}

	ctx := c.Request().Context()
	detail, err := r.srv.GetDetail(ctx, song)
	if err != nil {
		return r.e.Map(err)
	}

	return c.JSON(200, detail)
}

func (r *songRoutes) getSongs(c echo.Context) error {
	params := c.QueryParams()

	var group, song *string
	var releasedBefore, releasedAfter *time.Time

	if params.Has("group") {
		param := params.Get("group")
		group = &param
	}

	if params.Has("song") {
		param := params.Get("song")
		song = &param
	}

	if params.Has("releasedBefore") {
		parsedBefore, err := time.Parse(time.DateOnly, params.Get("releasedBefore"))
		if err != nil {
			return ErrBadQueryTime
		}
		releasedBefore = &parsedBefore
	}

	if params.Has("releasedAfter") {
		parsedAfter, err := time.Parse(time.DateOnly, params.Get("releasedAfter"))
		if err != nil {
			return ErrBadQueryTime
		}
		releasedAfter = &parsedAfter
	}

	filter := models.Filter{
		Group:          group,
		Song:           song,
		ReleasedBefore: releasedBefore,
		ReleasedAfter:  releasedAfter,
	}

	limit, err := strconv.Atoi(params.Get("limit"))
	if err != nil {
		return ErrBadQueryPagination
	}

	offset, err := strconv.Atoi(params.Get("offset"))
	if err != nil {
		return ErrBadQueryPagination
	}

	ctx := c.Request().Context()
	songs, err := r.srv.GetSongs(ctx, limit, offset, filter)
	if err != nil {
		return r.e.Map(err)
	}

	return c.JSON(200, songs)
}

func (r *songRoutes) getText(c echo.Context) error {
	params := c.QueryParams()

	if !params.Has("song") || !params.Has("group") || !params.Has("limit") || !params.Has("offset") {
		return ErrBadQuery
	}

	song := models.Song{
		GroupName: params.Get("group"),
		SongName:  params.Get("song"),
	}

	var limit, offset int

	limit, err := strconv.Atoi(params.Get("limit"))
	if err != nil {
		return ErrBadQueryPagination
	}

	offset, err = strconv.Atoi(params.Get("offset"))
	if err != nil {
		return ErrBadQueryPagination
	}

	ctx := c.Request().Context()
	text, err := r.srv.GetText(ctx, limit, offset, song)
	if err != nil {
		return r.e.Map(err)
	}

	return c.JSON(200, text)
}

func (r *songRoutes) deleteSong(c echo.Context) error {
	params := c.QueryParams()

	if !params.Has("song") || !params.Has("group") {
		return ErrBadQuery
	}

	song := models.Song{
		SongName:  params.Get("song"),
		GroupName: params.Get("group"),
	}

	ctx := c.Request().Context()
	if err := r.srv.Delete(ctx, song); err != nil {
		return r.e.Map(err)
	}

	return c.JSON(200, "Success!")
}

func (r *songRoutes) uploadSong(c echo.Context) error {
	params := c.QueryParams()

	if !params.Has("song") || !params.Has("group") || !params.Has("releaseDate") || !params.Has("link") || !params.Has("text") {
		return ErrBadQuery
	}

	releaseDate, err := time.Parse(time.DateOnly, params.Get("releaseDate"))
	if err != nil {
		return ErrBadQueryTime
	}

	song := models.SongWithDetailPlain{
		Song: models.Song{
			GroupName: params.Get("group"),
			SongName:  params.Get("song"),
		},
		SongDetail: models.SongDetail{
			ReleaseDate: releaseDate,
			Link:        params.Get("link"),
		},
		Text: params.Get("text"),
	}

	ctx := c.Request().Context()
	if err := r.srv.Create(ctx, song); err != nil {
		return r.e.Map(err)
	}

	return c.JSON(200, "Success!")
}

func (r *songRoutes) updateSong(c echo.Context) error {
	queryParams := c.QueryParams()
	formParams := c.Request().Form

	if !queryParams.Has("group") || !queryParams.Has("song") {
		return ErrBadQuery
	}

	if !formParams.Has("group") || !formParams.Has("song") {
		return ErrBadQuery
	}

	oldSong := models.Song{
		GroupName: queryParams.Get("group"),
		SongName:  queryParams.Get("song"),
	}

	parsedDate, err := time.Parse(time.DateOnly, formParams.Get("releaseDate"))
	if err != nil {
		return ErrBadQueryTime
	}

	newSong := models.SongWithDetailPlain{
		Song: models.Song{
			GroupName: formParams.Get("group"),
			SongName:  formParams.Get("song"),
		},
		SongDetail: models.SongDetail{
			ReleaseDate: parsedDate,
			Link:        formParams.Get("link"),
		},
		Text: formParams.Get("text"),
	}

	ctx := c.Request().Context()
	r.srv.Update(ctx, oldSong, newSong)

	return c.JSON(200, "Success!")

}
