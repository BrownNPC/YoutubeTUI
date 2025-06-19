package yt

import (
	"bytes"
	"encoding/json"
	"errors"
)

func GetPlaylist(playlistID string) (List, error) {
	const ytPlaylistUrl = "https://www.youtube.com/playlist?list="
	// try loading from cache
	if list, ok := loadFromCache(playlistID); ok {
		return list, nil
	}
	stdout, stderr, err := runYtDLP(
		"--flat-playlist", "--dump-single-json", ytPlaylistUrl+playlistID,
	)
	if stderr.Len() != 0 { // ytdlp error
		return List{}, errors.New(stderr.String())
	}
	if err != nil {
		return List{}, errors.Join(errors.New("failed to fetch playlist"), err)
	}
	var p List
	err = json.Unmarshal(stdout.Bytes(), &p)
	if err != nil {
		return p, err
	}
	saveToCache(p)
	return p, nil
}
func GetStreamURL(videoURL string) (url string, err error) {
	// yt-dlp -f "bestaudio[ext=webm][acodec=opus]" -g
	var stdout, stderr bytes.Buffer
	stdout, stderr, err = runYtDLP(
		"-f", "bestaudio[ext=webm][acodec=opus]", "-g", videoURL)
	if stderr.Len() != 0 { // ytdlp error
		err = errors.New(stderr.String())
		return
	}
	if err != nil {
		return
	}
	runes := []rune(stdout.String())
	// the last character is a newline and that really messes things up
	if runes[len(runes)-1] == '\n' {
		runes = runes[:len(runes)-1]
	}
	return string(runes), nil
}
