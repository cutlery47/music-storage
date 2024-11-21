package app

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cutlery47/music-storage/internal/models"
	"github.com/cutlery47/music-storage/internal/repository"
	"github.com/cutlery47/music-storage/internal/service"
)

func Run() error {
	ctx := context.Background()

	url := "postgresql://postgres:postgres@localhost:5433/music?sslmode=disable"

	repo, err := repository.NewMusicRepository(url)
	if err != nil {
		return fmt.Errorf("couldn't connect to the database: %v", err)
	}

	srv := service.NewMusicService(repo)

	// date := time.Date(2001, 10, 8, 0, 0, 0, 0, time.UTC)
	// group := "Linkin Park"

	filter := models.Filter{}

	songs, err := srv.GetSongs(ctx, 10, 0, filter)
	if err != nil {
		return fmt.Errorf("srv.GetSongs: %v", err)
	}

	log.Println("songs:", songs)

	song := models.Song{
		GroupName: "Coldplay",
		SongName:  "Fix You",
	}

	info, err := srv.GetDetail(ctx, song)
	if err != nil {
		return fmt.Errorf("srv.GetDetail: %v", err)
	}

	log.Println("info:", info)

	text, err := srv.GetText(ctx, 1, 1, song)
	if err != nil {
		return fmt.Errorf("srv.GetText: %v", err)
	}

	log.Println("text:", text)

	newSong := models.SongWithDetailPlain{
		Song: models.Song{
			GroupName: "Kendrick Lamar",
			SongName:  "6:16 in LA",
		},
		SongDetail: models.SongDetail{
			ReleaseDate: time.Now(),
			Link:        "https://example.com/xyu",
		},
		Text: "I fuck Drake aye, aye\n I love kids aye, aye\n",
	}

	if err := srv.Create(ctx, newSong); err != nil {
		return fmt.Errorf("srv.Create: %v", err)
	}

	newText, err := srv.GetText(ctx, 10, 0, newSong.Song)
	if err != nil {
		return fmt.Errorf("srv.GetText: %v", err)
	}

	fmt.Println("text:", newText)

	updSong := models.SongWithDetailPlain{
		Song: models.Song{
			GroupName: "Kendrick Lamar",
			SongName:  "6:17 in LA",
		},
		SongDetail: models.SongDetail{
			ReleaseDate: time.Now(),
			Link:        "https://example.com/xyu",
		},
		Text: "I fuck didi\n I love drake\n I fuck didi\n I love drake\n I fuck didi\n",
	}

	if err := srv.Update(ctx, newSong.Song, updSong); err != nil {
		return fmt.Errorf("srv.Update: %v", err)
	}

	anotherText, err := srv.GetText(ctx, 2, 4, updSong.Song)
	if err != nil {
		return fmt.Errorf("srv.GetText: %v", err)
	}

	fmt.Println("another text:", anotherText)

	if err := srv.Delete(ctx, updSong.Song); err != nil {
		return fmt.Errorf("srv.Delete: %v", err)
	}

	return nil
}
