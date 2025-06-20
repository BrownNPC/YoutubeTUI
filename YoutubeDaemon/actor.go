package daemon

import (
	"fmt"
	"net/http"
	"sync"
	"ytt/YoutubeDaemon/yt"

	"github.com/ebitengine/oto/v3"
	"github.com/jfbus/httprs"
)

type Event = any

type EventTrackStarted Track
type EventPlaylistStarted Playlist
type EventErr = error
type EventInfo = string

type Command any
type CmdStop struct{}
type CmdFetchStreamURL struct {
	track *Track
}
type CmdPlayTrack struct {
	track Track
}
type CmdRegisterPlaylists struct {
	playlistIDs []string
}
type CmdGetRegisteredPlaylists struct {
	playlists chan<- []Playlist
}

var cmdCh chan Command
var events chan Event

func InitDaemon() {
	<-yt.Ready
	cmdCh = make(chan Command)
	events = make(chan Event, 1)
	go playerManager(cmdCh, events)
}
func Events() <-chan Event {
	return events
}

func playerManager(cmdCh <-chan Command, events chan<- Event) {
	var (
		player    *oto.Player
		cleanup   func()
		playlists = []Playlist{} // registered playlists
	)

	for cmd := range cmdCh {
		switch cmd := cmd.(type) {
		case CmdStop:
			if cleanup != nil {
				cleanup()
				cleanup = nil
				events <- fmt.Sprintln("[INFO] asked to stop")
			}
		case CmdFetchStreamURL:
			var t *Track = cmd.track
			events <- fmt.Sprintf("[INFO] Trying to fetch stream url for %v\n", *t)
			if t.StreamingURL != "" {
				continue
			}
			url, err := yt.GetStreamURL(t.VideoURL)
			if err != nil {
				events <- err
				continue
			}
			events <- fmt.Sprintf("[INFO] Fetched streaming URL %s\n", url)
			t.StreamingURL = url
		case CmdPlayTrack:
			var t Track = cmd.track
			if t.StreamingURL == "" {
				events <- fmt.Errorf("Trying to play but streaming url is empty for %v\n", t)
				continue
			}
			events <- fmt.Sprintf("[INFO] Getting response body for track %v\n", t)
			resp, err := http.Get(t.StreamingURL)
			if err != nil {
				events <- err
				continue
			}
			f := httprs.NewHttpReadSeeker(resp)
			reader, _, err := newWebMReader(f)
			events <- fmt.Sprintf("[INFO] Decoder initialized for %v\n", t)
			if err != nil {
				events <- err
				continue
			}
			player = otoCtx.NewPlayer(reader)
			player.Play()
			events <- fmt.Sprintf("[INFO] player is playing for %v\n", t)
			if cleanup != nil {
				panic("assert: cleanup should be nil before playing track")
			}
			cleanup = func() {
				player.Close()
				f.Close()
				reader.Close()
			}
		case CmdGetRegisteredPlaylists:
			cmd.playlists <- playlists
		case CmdRegisterPlaylists:
			added := make(chan Playlist, 100)
			semaphore := make(chan struct{}, 3) //limit to 3 playlists being fetched
			var wg sync.WaitGroup
			for _, id := range cmd.playlistIDs {
				wg.Add(1)
				go func() {
					defer wg.Done()
					semaphore <- struct{}{}
					defer func() { <-semaphore }() // release

					list, err := yt.GetPlaylist(id)
					if err != nil {
						err = fmt.Errorf("fetching playlist %s: %w", id, err)
						events <- err
						return
					}
					events <- err
					pl := Playlist{List: list}
					for _, t := range list.Entries {
						pl.Tracks = append(pl.Tracks, Track{Entry: t})
					}
					added <- pl
				}()
			}
			wg.Wait()
			close(added)
			for p := range added {
				playlists = append(playlists, p)
			}
		}
	}
}
func RegisterPlaylists(playlistIDs ...string) {
	cmdCh <- CmdRegisterPlaylists{playlistIDs}
}
func GetRegisteredPlaylists() []Playlist {
	playlistsCh := make(chan []Playlist)
	cmdCh <- CmdGetRegisteredPlaylists{playlistsCh}
	return <-playlistsCh
}
func PlayTrack(t *Track) {
	cmdCh <- CmdFetchStreamURL{t}
	// play the track
	cmdCh <- CmdStop{}
	cmdCh <- CmdPlayTrack{track: *t}
}
