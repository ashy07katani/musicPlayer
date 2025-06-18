package util

import (
	"encoding/json"
	"fmt"
	"io"
	"music-player/model"
	"net/http"
	"net/url"
	"os/exec"
	"strings"
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

func ExtractAlbumArt(albumName string) (string, error) {
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
	if err != nil {
		return "", err
	}

	var coverArtResponse model.CoverArtArchieveResponse
	if err := json.Unmarshal(body, &coverArtResponse); err != nil {
		return "", err
	}

	if len(coverArtResponse.Images) > 0 {
		return coverArtResponse.Images[0].Thumbnails.Mid, nil
	}

	return "", fmt.Errorf("no album art found for album: %s", albumName)
}
