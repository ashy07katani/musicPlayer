package model

import (
	"fmt"
	"strconv"
	"time"
)

type SongMetaData struct {
	Format *Format `json:"format"`
}

type Format struct {
	Duration string `json:"duration"`
	Tags     *Tags  `json:"tags"`
}

type Tags struct {
	Title        string `json:"title"`
	Artist       string `json:"artist"`
	Album        string `json:"album"`
	Genre        string `json:"genre"`
	PlaylistPath string `json:"playlistPath"`
	AlbumArt     string `json:"albumArt"`
}

type MusicbrainzResponse struct {
	Releases []*Release `json:"releases"`
}

type Release struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

type CoverArtArchieveResponse struct {
	Images []*Image `json:"images"`
}

type Image struct {
	Thumbnails *Thumbnail `json:"thumbnails"`
}
type Thumbnail struct {
	Mid   string `json:"500"`
	Small string `json:"250"`
}

type Song struct {
	ID           uint
	Title        string
	Artist       string
	Album        string
	Genre        string
	Duration     float64
	PlaylistPath string
	AlbumArt     string
	CreatedAt    time.Time
}

func NewSong(songMetaData *SongMetaData) *Song {
	song := new(Song)
	song.Title = songMetaData.Format.Tags.Title
	song.Artist = songMetaData.Format.Tags.Artist
	song.Album = songMetaData.Format.Tags.Album
	song.Genre = songMetaData.Format.Tags.Genre
	song.AlbumArt = songMetaData.Format.Tags.AlbumArt
	duration, err := strconv.ParseFloat(songMetaData.Format.Duration, 64)
	song.Duration = duration
	fmt.Println(song.Artist, songMetaData.Format.Tags.Artist)
	if err != nil {
		panic("Error converting string to float")
	}
	song.PlaylistPath = songMetaData.Format.Tags.PlaylistPath
	song.CreatedAt = time.Now()
	return song
	//songMetaData.Format.Duration
}
