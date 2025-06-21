package daemon

import (
	"fmt"
	"net/http"
	"slices"
	"sync"
	"time"
	"ytt/YoutubeDaemon/yt"

	"github.com/ebitengine/oto/v3"
	"github.com/jfbus/httprs"
)

type Event any

type EventTrackStarted Track
type EventErr = error
type EventInfo = string

type Command any
type CmdStop struct{}
type CmdPlayNextTrack struct{}
type CmdFetchStreamURL struct{ *Track }
type CmdPlayTrack struct{ *Track }
type CmdSetQueue struct{ Tracks []*Track }
type CmdGetQueue struct{ queue chan<- []*Track }
type CmdStartQueue struct{}               // start playing queue
type CmdSetQueuePosition struct{ *Track } // set queue to start from here
type CmdRegisterPlaylists struct{ playlistIDs []string }
type CmdGetRegisteredPlaylists struct{ playlists chan<- []Playlist }
type CmdGetCurrentTrackDuration struct{ duration chan<- time.Duration }

var cmdCh chan Command
var events chan Event

func InitDaemon() {
	<-yt.Ready
	cmdCh = make(chan Command)
	events = make(chan Event)
	go playerManager(cmdCh, events)
}
func Events() <-chan Event {
	return events
}

func playerManager(cmdCh <-chan Command, events chan<- Event) {
	var (
		player  *oto.Player
		cleanup func() //for stopping player

		playlists = []Playlist{} // registered playlists

		queue      = []*Track{} // []Track from within a playlist
		queueIndex int

		trackPlaying *Track
	)

	for cmd := range cmdCh {
		switch cmd := cmd.(type) {
		case CmdStop:
			if cleanup != nil {
				cleanup()
				cleanup = nil
				events <- fmt.Sprintln("[INFO] asked to stop")
			}
		case CmdSetQueue:
			queue = cmd.Tracks
			queueIndex = 0
		case CmdStartQueue:
			if len(queue) >= 1 {
				trackPlaying = queue[queueIndex]
			} else {
				events <- fmt.Errorf("Queue too small to play %d", len(queue))
				return
			}
			events <- fmt.Sprintf("[INFO] playing track %d: %s from queue", queueIndex, trackPlaying.Title)
			// TODO: Take cancelable context and pass it
			go PlayTrack(trackPlaying)
			queueIndex++
			queueIndex %= len(queue)
		case CmdSetQueuePosition:
			i := slices.Index(queue, cmd.Track)
			if i != -1 {
				queueIndex = i
			}
		case CmdFetchStreamURL:
			var t *Track = cmd.Track
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
			var t *Track = cmd.Track
			if t.StreamingURL == "" {
				events <- fmt.Errorf("Trying to play but streaming url is empty for %v\n", t)
				continue
			}
			events <- fmt.Sprintf("[INFO] Getting response body for track %s\n", t.Title)
			resp, err := http.Get(t.StreamingURL)
			if err != nil {
				events <- err
				continue
			}
			f := httprs.NewHttpReadSeeker(resp)
			reader, _, err := newWebMReader(f)
			events <- fmt.Sprintf("[INFO] Decoder initialized for %s\n", t.Title)
			if err != nil {
				events <- err
				continue
			}
			player = otoCtx.NewPlayer(reader)
			player.Play()
			events <- fmt.Sprintf("[INFO] player is playing %s\n", t.Title)
			events <- EventTrackStarted(*t)
			trackPlaying = t
			if cleanup != nil {
				panic("assert: cleanup should be nil before playing track")
			}
			cleanup = func() {
				player.Close()
				f.Close()
				reader.Close()
			}
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
						if err != nil {
							events <- err
						}
						return
					}
					events <- err
					pl := Playlist{List: list}
					for _, t := range list.Entries {
						pl.Tracks = append(pl.Tracks, &Track{Entry: t})
					}
					added <- pl
				}()
			}
			wg.Wait()
			close(added)
			for p := range added {
				playlists = append(playlists, p)
			}
		case CmdGetQueue:
			cmd.queue <- queue
		case CmdGetRegisteredPlaylists:
			cmd.playlists <- playlists
		case CmdGetCurrentTrackDuration:
			cmd.duration <- time.Second *
				time.Duration(trackPlaying.DurationSeconds)
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

// TODO: make this take a ctx that is cancelable
// or an "ok" channel. The goal is to show a loading screen while
// we prepare the player, and show a cancel button to the user
func PlayTrack(t *Track) {
	cmdCh <- CmdStop{}
	cmdCh <- CmdFetchStreamURL{t}
	// play the track
	cmdCh <- CmdPlayTrack{t}
}

// TODO: make this take a ctx that is cancelable
// or an "ok" channel. The goal is to show a loading screen while
// we prepare the player, and show a cancel button to the user
func PlayPlaylist(p Playlist) {
	cmdCh <- CmdStop{}
	cmdCh <- CmdSetQueue{p.Tracks}
	cmdCh <- CmdStartQueue{}
}

func GetQueue() []*Track {
	queue := make(chan []*Track)
	cmdCh <- CmdGetQueue{queue}
	return <-queue
}
