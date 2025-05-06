package daemon

import (
	"encoding/binary"
	"errors"
	"io"
	"time"

	"github.com/ebml-go/webm"
	"github.com/hraban/opus"
)

// Reader encapsulates the audio decoding logic and implements io.Reader.
type Reader struct {
	pr         *io.PipeReader // Pipe reader for audio data
	webmReader *webm.Reader
	webmFile   webm.WebM
	Progress   time.Duration
	quit       bool
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

	decoder, err := opus.NewDecoder(int(track.SamplingFrequency), int(track.Channels))
	if err != nil {
		return nil, nil, err
	}

	decodeBuffer := make([]float32, 1000*int(track.Channels))

	r := &Reader{
		pr:         pr,
		webmReader: webmReader,
		webmFile:   webmFile,
	}
	go r.decode(pw, webmReader, decoder, decodeBuffer, track)

	return r, track, nil
}

func (r *Reader) decode(pw *io.PipeWriter, webmReader *webm.Reader, decoder *opus.Decoder, decodeBuffer []float32, track *webm.TrackEntry) {
	defer pw.Close()
	for {
		packet := <-webmReader.Chan
		if r.quit { // reader is closed
			break
		}
		r.Progress = packet.Timecode

		nSamples, err := decoder.DecodeFloat32(packet.Data, decodeBuffer)
		if nSamples == 0 { //important or audio will stop playing on seek
			continue
		}
		if err != nil {
			pw.CloseWithError(err)
		}

		// Convert float32 samples to bytes and write to the pipe
		err = binary.Write(pw, binary.LittleEndian, decodeBuffer[:nSamples*int(track.Channels)])
		if err != nil {
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
