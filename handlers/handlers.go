package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"music-player/model"
	"music-player/repo"
	"music-player/util"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

type MusicHandler struct {
	DB          *sql.DB
	RedisClient *redis.Client
}

func NewMusicHandler(db *sql.DB, redisClient *redis.Client) *MusicHandler {
	handler := new(MusicHandler)
	handler.DB = db
	handler.RedisClient = redisClient
	return handler
}

func (h *MusicHandler) StreamMusic(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fileName := params["filename"]
	filePath := "C:/Users/tripa/OneDrive/Documents/MusicPlayer/" + fileName
	file, err := os.Open(filePath)
	if err != nil {
		log.Panicln("Error opening file ", fileName)
	}
	defer file.Close()
	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()
	Range := r.Header.Get("Range")
	if Range == "" {
		w.Header().Set("Content-Type", "audio/mpeg")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", fileSize))
		http.ServeContent(w, r, fileName, fileInfo.ModTime(), file)
		return
	}

	var start, end int64
	if strings.HasPrefix(Range, "bytes=") {
		ranges := strings.Split(strings.TrimPrefix(Range, "bytes="), "-")
		start, _ = strconv.ParseInt(ranges[0], 10, 64)
		if ranges[1] != "" {
			end, _ = strconv.ParseInt(ranges[1], 10, 64)
		}
	}
	if end == 0 || end >= fileSize {
		end = fileSize - 1
	}
	chunkSize := end - start + 1

	file.Seek(start, io.SeekStart)

	w.Header().Set("Content-Type", "audio/mpeg")
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", chunkSize))
	w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	w.WriteHeader(http.StatusPartialContent)

	io.CopyN(w, file, chunkSize)
	fmt.Println(fileName)
}

func (h *MusicHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	file, fileheader, err := r.FormFile("musicFile")
	//create a temporary upload folder
	defer file.Close()
	os.MkdirAll("uploads", os.ModePerm)
	savedFilePath := "uploads/" + fileheader.Filename
	dst, err := os.Create(savedFilePath)
	if err != nil {
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	playlistPath := ""
	go func() {
		segmentTime := "5"
		outputDir := "chunks"
		os.MkdirAll(outputDir, os.ModePerm)
		log.Println(fileheader.Filename, fileheader.Header, fileheader.Size)
		filename := strings.TrimSuffix(fileheader.Filename, ".mp3")
		segmentPattern := fmt.Sprintf("%s/%s_%%04d.ts", outputDir, filename)
		playlistPath = fmt.Sprintf("%s/%s_playlist.m3u8", outputDir, filename)
		cmd := exec.Command(
			"ffmpeg",
			"-i", savedFilePath,
			"-acodec", "aac",
			"-vn",
			"-f", "hls",
			"-hls_time", segmentTime,
			"-hls_list_size", "0",
			"-hls_segment_filename", segmentPattern,
			"-hls_base_url", "/chunks/",
			playlistPath,
		)
		err = cmd.Run()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}()
	info, err := util.ExtractMetadata(savedFilePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	albumArtUrl, err := util.ExtractAlbumArt(info.Format.Tags.Artist, info.Format.Tags.Album, h.RedisClient)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if albumArtUrl != "" {
		info.Format.Tags.AlbumArt = albumArtUrl
	}
	info.Format.Tags.PlaylistPath = playlistPath
	songEntry := model.NewSong(info)
	fmt.Println(songEntry.Artist)
	err = repo.InsertSongMetaData(h.DB, songEntry)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func (h *MusicHandler) StreamHLS(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]

	// Defensive: If someone passes a chunk filename like DawnFM_0000.ts, reject it
	if strings.Contains(filename, ".ts") || strings.Contains(filename, "_0000") {
		http.Error(w, "Invalid playlist request", http.StatusBadRequest)
		return
	}

	playlistPath := fmt.Sprintf("chunks/%s_playlist.m3u8", filename)

	// Check if file exists
	if _, err := os.Stat(playlistPath); os.IsNotExist(err) {
		fmt.Println("Playlist not found:", playlistPath)
		http.Error(w, "Playlist not found", http.StatusNotFound)
		return
	}

	// Set appropriate headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	http.ServeFile(w, r, playlistPath)
}
