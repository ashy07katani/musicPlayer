package repo

import (
	"database/sql"
	"fmt"
	"log"
	"music-player/model"
)

func InsertSongMetaData(db *sql.DB, song *model.Song) error {
	fmt.Println("song albumart is ", song.AlbumArt)
	_, err := db.Exec(INSERT_SONG, song.Title, song.Artist, song.Album, song.Genre, song.Duration, song.PlaylistPath, song.AlbumArt, song.CreatedAt)
	if err != nil {
		log.Println("Insert failed:", err)
	}
	return err
}

//INSERT_SONG = `INSERT INTO songsInfo VALUES (title, artist, album, genre, duration,  playlist_path, album_art, created_at ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
