package models

import "time"

type Song struct {
	GroupName string `db:"group_name"`
	SongName  string `db:"song_name"`
}

type SongDetail struct {
	ReleaseDate time.Time `db:"released_at"`
	Link        string    `db:"link"`
}

// необработанный текст песни
type SongPlainText string

// обработанный текст песни
type SongSplitText []string

type SongWithDetail struct {
	Song
	SongDetail
}

type SongWithDetailPlain struct {
	Song
	SongDetail
	SongPlainText
}

type SongWithDetailSplit struct {
	Song
	SongDetail
	SongSplitText
}
