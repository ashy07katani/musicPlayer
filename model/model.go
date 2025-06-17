package model

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
