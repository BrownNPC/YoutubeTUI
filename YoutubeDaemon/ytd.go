package daemon

import (
	"log"
	"ytt/YoutubeDaemon/yt"

	"github.com/ebitengine/oto/v3"
)

var (
	otoCtx *oto.Context
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

func init() {
	op := oto.NewContextOptions{
		SampleRate:   48000,
		ChannelCount: 2,
		Format:       oto.FormatFloat32LE,
	}

	ctx, ready, err := oto.NewContext(&op)
	if err != nil {
		log.Fatal("could not initialize audio", err)
	}
	<-ready
	otoCtx=ctx
}
