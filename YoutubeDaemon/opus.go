package daemon

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"time"
	"ytt/YoutubeDaemon/opus"

	"github.com/ebml-go/webm"
)

// Reader encapsulates the audio decoding logic and implements io.Reader.
type Reader struct {
	pr         *io.PipeReader // Pipe reader for audio data
	webmReader *webm.Reader
	webmFile   webm.WebM
	Progress   time.Duration
	quit       bool
}

var decoder opus.Decoder

func init() {
	dec, err := opus.NewDecoder(48000, 2)
	if err != nil {
		panic(fmt.Sprintln("failed to create decoder", err))
	}
	decoder = dec
}

// newWebMReader initializes a new Reader by parsing the WebM file and starting a decoding goroutine.
func newWebMReader(rs io.ReadSeeker) (*Reader, *webm.TrackEntry, error) {
	var webmFile webm.WebM
	webmReader, err := webm.Parse(rs, &webmFile)

	if err != nil {
		return nil, nil, err
	}
	track := webmFile.FindFirstAudioTrack()
	if track == nil {
		return nil, nil, errors.New("no audio track found")
	}

	pr, pw := io.Pipe()

	decodeBuffer := make([]float32, 1000*int(track.Channels))

	r := &Reader{
		pr:         pr,
		webmReader: webmReader,
		webmFile:   webmFile,
	}
	go r.decode(pw, webmReader, decodeBuffer, track)

	return r, track, nil
}

func (r *Reader) decode(pw *io.PipeWriter, webmReader *webm.Reader, decodeBuffer []float32, track *webm.TrackEntry) {
	defer pw.Close()
	for {
		packet := <-webmReader.Chan
		if r.quit { // reader is closed
			break
		}
		r.Progress = packet.Timecode
		events <- fmt.Sprintln(packet.Timecode.Seconds())
		nSamples, err := decoder.DecodeFloat32(packet.Data, decodeBuffer)
		if nSamples == 0 { //important or audio will stop playing on seek
			continue
		}
		if err != nil {
			events <- err
			pw.CloseWithError(err)
		}

		// Convert float32 samples to bytes and write to the pipe
		err = binary.Write(pw, binary.LittleEndian, decodeBuffer[:nSamples*int(track.Channels)])
		if err != nil {
			events <- err
			pw.CloseWithError(err)
		}
	}
}

// Read implements the io.Reader interface by reading from the pipe.
func (r *Reader) Read(data []byte) (int, error) {
	return r.pr.Read(data)
}

func (r *Reader) Seek(t time.Duration) {
	r.webmReader.Seek(t)
}
func (r *Reader) Close() {
	r.quit = true
}
