package util

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"music-player/model"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

func ExtractMetadata(filePath string) (*model.SongMetaData, error) {
	cmd := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", filePath)
	out, err := cmd.Output()
	fmt.Println(string(out))
	if err != nil {
		return nil, err
	}
	meta := new(model.SongMetaData)
	if err := json.Unmarshal(out, meta); err != nil {
		return nil, err
	}
	return meta, nil
}

func ExtractAlbumArt(artistName string, albumName string, redisClient *redis.Client) (string, error) {

	redisKey := normalise(albumName, artistName)
	ctx := context.Background()
	val, err := redisClient.Get(ctx, redisKey).Result()
	if err == nil {
		return val, err
	}
	if err == redis.Nil {
		client := &http.Client{}

		// Escape the album name to ensure valid URL
		escapedAlbumName := url.QueryEscape(albumName)
		musicbrainzURL := fmt.Sprintf("https://musicbrainz.org/ws/2/release/?query=release:%s&fmt=json", escapedAlbumName)

		// Step 1: MusicBrainz request
		req, err := http.NewRequest(http.MethodGet, musicbrainzURL, nil)
		if err != nil {
			return "", err
		}
		req.Header.Set("User-Agent", "MySpotifyClone/0.1 (ashykatani@gmail.com)")

		resp, err := client.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		// Check response type
		contentType := resp.Header.Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			bodyBytes, _ := io.ReadAll(resp.Body)
			fmt.Println("Non-JSON response:\n", string(bodyBytes))
			return "", fmt.Errorf("expected JSON, got: %s", contentType)
		}

		body, err := io.ReadAll(resp.Body)
		//fmt.Println(string(body))
		if err != nil {
			return "", err
		}

		var musicBrainzResponse model.MusicbrainzResponse
		if err := json.Unmarshal(body, &musicBrainzResponse); err != nil {
			return "", err
		}

		if len(musicBrainzResponse.Releases) == 0 {

			return "", fmt.Errorf("no releases found for album: %s", albumName)
		}

		releaseID := musicBrainzResponse.Releases[0].Id
		coverArtURL := fmt.Sprintf("https://coverartarchive.org/release/%s", releaseID)

		// Step 2: Cover Art Archive request
		req, err = http.NewRequest(http.MethodGet, coverArtURL, nil)
		if err != nil {
			return "", err
		}
		req.Header.Set("User-Agent", "MySpotifyClone/0.1 (ashykatani@gmail.com)")

		resp, err = client.Do(req)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		body, err = io.ReadAll(resp.Body)
		//fmt.Println(string(body))
		if err != nil {
			return "", err
		}

		var coverArtResponse model.CoverArtArchieveResponse
		if err := json.Unmarshal(body, &coverArtResponse); err != nil {
			return "", err
		}

		if len(coverArtResponse.Images) > 0 {
			coverArt := coverArtResponse.Images[0].Thumbnails.Mid
			err = redisClient.Set(context.Background(), normalise(albumName, artistName), coverArt, time.Second*3600*24).Err()
			if err != nil {
				fmt.Println("couldn't store the value in redis")
			}
			return coverArt, nil
		}

	}
	return "", fmt.Errorf("no album art found for album: %s", albumName)

}

func normalise(album string, artist string) string {
	key := fmt.Sprintf("%s-%s", artist, album)
	key = strings.ToLower(key)
	key = strings.ReplaceAll(key, " ", "")
	return key
}
