package main

import (
	"time"
	"ytt/daemon"

	tea "github.com/charmbracelet/bubbletea"
)

type TickMsg struct{}

func CmdTick() tea.Msg {
	<-time.Tick(time.Millisecond * 250)
	return TickMsg{}
}

// seek forward or backward
func CmdSeek(forward bool) {
	seekTime := time.Second * 10
	if forward {
		daemon.Seek(seekTime)
	} else {
		daemon.Seek(seekTime * -1)
	}
}
