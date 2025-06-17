package util

import (
	"encoding/json"
	"fmt"
	"music-player/model"
	"os/exec"
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
