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
