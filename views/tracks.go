package views

import (
	daemon "ytt/YoutubeDaemon"
	"ytt/components"
	"ytt/themes"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

// does not play anything, it just allows you to add tracks to queue
type TracksModel struct {
	initialized   bool
	width, height int

	list components.List
}

var ActiveTracksModel = TracksModel{}
var SelectedTrack *daemon.Track

func (m *TracksModel) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.list, _ = m.list.Update(msg)
	case tea.KeyMsg:
		if e, ok := m.list.Hovered(); ok {
			if msg.String() == "enter" {
				SelectedTrack, ok = e.CustomData.(*daemon.Track)
				if !ok {
					panic("assert: custom data of Track was not found")
				}
				go daemon.PlayTrack(SelectedTrack)
			}
		}
	case tea.MouseMsg:
		if e, ok := m.list.MouseHovered(msg); ok {
			if msg.Mouse().Button == tea.MouseLeft {
				SelectedTrack, ok = e.CustomData.(*daemon.Track)
				if !ok {
					panic("assert: custom data of Track was not found")
				}
				daemon.PlayTrack(SelectedTrack)
			}
		}
	}
	if m.initialized {
		m.list, cmd = m.list.Update(msg)
	}
	return cmd
}
func (m TracksModel) View() string {
	t := themes.Active()
	var base = lipgloss.NewStyle().
		Background(t.Background)
	var o string
	if !m.initialized {
		return "Please select a playlist first"
	}
	o += m.list.View()
	return base.
		Width(m.width).
		Height(m.height).
		Render(o)
}
