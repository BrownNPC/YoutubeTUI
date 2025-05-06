package daemon

import (
	"log"
	"net/http"

	"github.com/ebitengine/oto/v3"
	"github.com/jfbus/httprs"
)

func newPlayerFromUrl(url string) (func(), *oto.Player, *Reader) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	f := httprs.NewHttpReadSeeker(resp)

	reader, _, err := newWebMReader(f)
	if err != nil {
		log.Fatal(err)
	}

	player := otoCtx.NewPlayer(reader)
	player.Play()
	var Close = func() {
		f.Close()
		player.Close()
		reader.Close()
	}
	return Close, player, reader
}
