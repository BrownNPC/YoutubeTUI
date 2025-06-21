package yt

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
)

var ytdlpPath string

const (
	xdgCacheDir = "ytt" // Cache directory that will be appended to the XDG cache directory.
)

var Ready = make(chan struct{})

func init() {
	go func() {
		install, err := Install(context.TODO(), nil)
		if err != nil {
			fmt.Println("Error installing ytdlp", err)
			os.Exit(1)
		}
		ytdlpPath = install.Executable
		Ready <- struct{}{}
	}()
}
func runYtDLP(args ...string) (stdoutBuf, stderrBuf bytes.Buffer, err error) {
	args = append(args, "--quiet", "--no-warnings") // only errors in stderr
	cmd := exec.Command(ytdlpPath, args...)
	// Return values
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	err = cmd.Run()
	return
}

// playlist or search results
type List struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Channel     string  `json:"channel"`
	Entries     []Entry `json:"entries"`
}

// track inside playlist or search results list
type Entry struct {
	ID              string `json:"id"`
	VideoURL        string `json:"url"`
	Title           string `json:"title"`
	DurationSeconds int    `json:"duration"`
	ChannelURL      string `json:"channel_url"`
	Uploader        string `json:"uploader"`
	ViewCount       int    `json:"view_count"`
}
