package main

import (
	"time"

	tea "github.com/charmbracelet/bubbletea/v2"
)

type TickMsg struct{}

func CmdTick() tea.Msg {
	<-time.Tick(time.Millisecond * 250)
	return TickMsg{}
}

