package daemon

import (
	"context"
	"sync"
	"ytt/YoutubeDaemon/yt"

	"github.com/ebitengine/oto/v3"
)

var (
	mu     sync.Mutex
	otoCtx *oto.Context
)

var Player = struct {
	plr     *oto.Player
	queue   []Track
	shuffle bool
	repeat  RepeatMode
	Events  chan Event
	close   context.CancelFunc
}{
	Events: make(chan Event, 100),
	queue:  make([]Track, 0),
}

type PlayerState string

// Repeat modes
type RepeatMode int

const (
	RepeatOff RepeatMode = iota
	RepeatOne
	RepeatAll
)
const (
	StateStopped PlayerState = "stopped"
	StatePlaying PlayerState = "playing"
	StatePaused  PlayerState = "paused"
)

// A playlist is just an ordered slice of Tracks
type Playlist struct {
	yt.List
	Tracks []Track
}
type Track struct {
	yt.Entry
	StreamingURL string
}

type Event = any

type EventTrackStarted Track
type EventPlaylistStarted Playlist
type EventPlaylistRegistered Playlist
type EventErr = error
