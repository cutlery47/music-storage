package v1

import (
	"fmt"
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

// @Summary 		Get Info
// @Description 	Get info about a particular song
// @Tags 			Songs
// @Param			group		query		string		true    "desired group"
// @Param 			song		query		string		true	"desired song"
// @Success			200 		{object} 	models.SongDetail
// @Failure 		400			{object}    echo.HTTPError
// @Failure			404			{object}	echo.HTTPError
// @Failure			500			{object} 	echo.HTTPError
// @Router 			/api/v1/songs/info [get]
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
		fmt.Println(err)
		return r.e.Map(err)
	}

	return c.JSON(200, detail)
}

// @Summary 		Get Songs
// @Description 	Get songs by specified filters
// @Tags 			Songs
// @Param			group				query		string		false   "desired group"
// @Param 			song				query		string		false   "desired song"
// @Param			releasedBefore		query		string		false   "upper time-bound for when the song was released"
// @Param 			releasedAfter		query		string		false	"lower time-bound for when the song was released"
// @Param			limit				query		int			true    "pagination limit"
// @Param 			offset				query		int			true	"pagination offset"
// @Success			200 				{object} 	[]models.SongWithDetail
// @Failure 		400					{object}    echo.HTTPError
// @Failure			404					{object}	echo.HTTPError
// @Failure			500					{object} 	echo.HTTPError
// @Router 			/api/v1/songs [get]
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

// @Summary 		Get Texts
// @Description 	Get specified songs' lyrics
// @Tags 			Songs
// @Param			group				query		string		true   "desired group"
// @Param 			song				query		string		true   "desired song"
// @Param			limit				query		int			true    "pagination limit"
// @Param 			offset				query		int			true	"pagination offset"
// @Success			200 				{object} 	string
// @Failure 		400					{object}    echo.HTTPError
// @Failure			404					{object}	echo.HTTPError
// @Failure			500					{object} 	echo.HTTPError
// @Router 			/api/v1/songs/text [get]
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

// @Summary 		Delete Song
// @Description 	Delete specific song
// @Tags 			Songs
// @Param			group				query		string		true   "desired group"
// @Param 			song				query		string		true   "desired song"
// @Success			200 				{object} 	string
// @Failure 		400					{object}    echo.HTTPError
// @Failure			404					{object}	echo.HTTPError
// @Failure			500					{object} 	echo.HTTPError
// @Router 			/api/v1/songs [delete]
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

// @Summary 		Upload Song
// @Description 	Upload a new song
// @Tags 			Songs
// @Param			group				formData	string		true    "desired group"
// @Param 			song				formData	string		true    "desired song"
// @Param			releaseDate			formData	string		true    "song release date"
// @Param 			link				formData	string		true	"link to some media"
// @Param 			text				formData	string		true	"song lyrics"
// @Success			200 				{object} 	string
// @Failure 		400					{object}    echo.HTTPError
// @Failure			404					{object}	echo.HTTPError
// @Failure			500					{object} 	echo.HTTPError
// @Router 			/api/v1/songs [post]
func (r *songRoutes) uploadSong(c echo.Context) error {
	params, _ := c.FormParams()

	if !params.Has("song") || !params.Has("group") || !params.Has("releaseDate") || !params.Has("link") || !params.Has("text") {
		return ErrBadBody
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

// @Summary 		Update Song
// @Description 	Upload specific song data
// @Tags 			Songs
// @Param			group				query		string		true    "initial group"
// @Param 			song				query		string		true    "initial song"
// @Param			group				formData	string		true    "edited group"
// @Param 			song				formData	string		true    "edited song"
// @Param			releaseDate			formData	string		true    "edited release date"
// @Param 			link				formData	string		true	"edited link to some media"
// @Param 			text				formData	string		true	"edited song lyrics"
// @Success			200 				{object} 	string
// @Failure 		400					{object}    echo.HTTPError
// @Failure			404					{object}	echo.HTTPError
// @Failure			500					{object} 	echo.HTTPError
// @Router 			/api/v1/songs [put]
func (r *songRoutes) updateSong(c echo.Context) error {
	queryParams := c.QueryParams()
	formParams, _ := c.FormParams()

	if !queryParams.Has("group") || !queryParams.Has("song") {
		return ErrBadQuery
	}

	if !formParams.Has("group") || !formParams.Has("song") {
		return ErrBadQuery
	}

	oldSong := models.Song{
		GroupName: queryParams["group"][0],
		SongName:  queryParams["song"][0],
	}

	parsedDate, err := time.Parse(time.DateOnly, formParams.Get("releaseDate"))
	if err != nil {
		return ErrBadQueryTime
	}

	newSong := models.SongWithDetailPlain{
		Song: models.Song{
			GroupName: formParams["group"][1],
			SongName:  formParams["song"][1],
		},
		SongDetail: models.SongDetail{
			ReleaseDate: parsedDate,
			Link:        formParams.Get("link"),
		},
		Text: formParams.Get("text"),
	}

	ctx := c.Request().Context()
	if err := r.srv.Update(ctx, oldSong, newSong); err != nil {
		return r.e.Map(err)
	}

	return c.JSON(200, "Success!")

}
