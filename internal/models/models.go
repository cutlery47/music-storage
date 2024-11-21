package models

import (
	"strings"
	"time"
)

type Song struct {
	GroupName string `db:"group_name"`
	SongName  string `db:"song_name"`
}

type SongDetail struct {
	ReleaseDate time.Time `db:"released_at"`
	Link        string    `db:"link"`
}

type SongWithDetail struct {
	Song
	SongDetail
}

// песня с необработанным текстом
type SongWithDetailPlain struct {
	Song
	SongDetail
	Text string
}

// песня с обработанным текстом
type SongWithDetailSplit struct {
	Song
	SongDetail
	Verses []string
}

// преобразование текста
// TODO: придумать более sophisticated подход
func (sp SongWithDetailPlain) Split() SongWithDetailSplit {
	splitText := strings.Split(sp.Text, "\n")

	return SongWithDetailSplit{
		Song:       sp.Song,
		SongDetail: sp.SongDetail,
		Verses:     splitText,
	}
}
