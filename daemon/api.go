// Package daemon is the audio player & the youtube music daemon
package daemon

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ebitengine/oto/v3"
)

var otoCtx *oto.Context
var currentPlayer *oto.Player
var mut sync.Mutex

// not using a struct because this only has one instance
var (
	currentTrackClose func()
	currentReader     *Reader

	seekDebounce      = time.Second / 4
	countseekdebounce time.Time

	Playlists       = map[string]Playlist{} // populated when we FetchPlaylist()
	currentPlaylist Playlist
)

// initialize audio player and download ytdlp
func init() {
	op := oto.NewContextOptions{
		SampleRate:   48000,
		ChannelCount: 2,
		Format:       oto.FormatFloat32LE,
	}
	var ready chan struct{}
	var err error
	otoCtx, ready, err = oto.NewContext(&op)
	if err != nil {
		log.Fatal("could not initialize audio", err)
	}
	<-ready
	err = DownloadYtdlp() // download ytdlp if not already downloaded
	if err != nil {
		panic(err)
	}
	countseekdebounce = time.Now()
}

func Track(url string) {
	currentTrackClose, currentPlayer, currentReader = newPlayerFromUrl(url)
}

// return time in 00:00:00 (hours:mins:seconds format)
func GetTimeStamp() string {
	mut.Lock()
	defer mut.Unlock()
	var t time.Duration
	if currentReader != nil {
		t = currentReader.Progress
	}
	hours := int(t.Hours())
	minutes := int(t.Minutes()) % 60
	seconds := int(t.Seconds()) % 60

	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func Seek(t time.Duration) {
	mut.Lock()
	defer mut.Unlock()
	if time.Since(countseekdebounce) < seekDebounce {
		return
	}
	countseekdebounce = time.Now()
	if currentReader != nil { // webm seeking is wonky
		if t > 0 {
			t /= 2
		} else if t < 0 {
			t *= 2
		}
		currentReader.Seek(currentReader.Progress + (t))
		currentReader.Progress += t
	}
}
func TogglePlayback() {
	mut.Lock()
	defer mut.Unlock()
	if currentPlayer == nil {
		return
	}
	if currentPlayer.IsPlaying() {
		currentPlayer.Pause()
	} else {
		currentPlayer.Play()
	}
}
func IsPlaying() bool {
	mut.Lock()
	defer mut.Unlock()
	if currentPlayer == nil {
		return false
	}
	return currentPlayer.IsPlaying()
}
func SetVolume(n float64) {
	mut.Lock()
	defer mut.Unlock()
	if currentPlayer == nil {
		return
	}
	if n > 150 {
		n = 150
	} else if n < 0 {
		n = 0
	}
	currentPlayer.SetVolume(n)
}
func GetVolume() float64 {
	mut.Lock()
	defer mut.Unlock()
	if currentPlayer == nil {
		return 0
	}
	return currentPlayer.Volume()
}
