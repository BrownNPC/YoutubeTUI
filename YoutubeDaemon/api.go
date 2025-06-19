package daemon

import (
	"context"
	"fmt"
	"sync"
	"time"
	"ytt/YoutubeDaemon/yt"
)

func PlayTrack(t *Track) {
	if t.StreamingURL == "" {
		var err error
		var url string

		maxRetries := 3
		for i := range maxRetries {
			url, err = yt.GetStreamURL(t.VideoURL)
			if err != nil {
				Player.Events <- EventErr(fmt.Errorf("failed to fetch stream URL (attempt %d): %w", i+1, err))
				time.Sleep(2 * time.Second) // Add delay between retries
				continue
			}
			t.StreamingURL = url
			break
		}

		if err != nil {
			return
		}
	}

	go beginStreaming(t.StreamingURL)
}

// only modified by addPlaylists
var registeredPlaylists = make([]Playlist, 0)

// get the playlists that have been added so far
func RegisteredPlaylists() []Playlist { return registeredPlaylists }
func AddPlaylists(playlistIds ...string) {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 3) //limit to 3 playlists being fetched

	for _, id := range playlistIds {
		wg.Add(1)
		go func() {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }() // release

			list, err := yt.GetPlaylist(id)
			if err != nil {
				Player.Events <- EventErr(fmt.Errorf("fetching playlist %s: %w", id, err))
				return
			}
			pl := Playlist{List: list}
			for _, t := range list.Entries {
				pl.Tracks = append(pl.Tracks, Track{Entry: t})
			}

			Player.Events <- EventPlaylistRegistered(pl)
			registeredPlaylists = append(registeredPlaylists, pl)
		}()
	}
	wg.Wait()
}
func Events() chan Event { return Player.Events }

func PlayPlaylist(p Playlist) {
	if Player.close != nil {
		Player.close()
	}
	mu.Lock()
	Player.queue = p.Tracks
	mu.Unlock()
	playQueue()
}
