package daemon

import (
	"log"

	"github.com/ebitengine/oto/v3"
)

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
