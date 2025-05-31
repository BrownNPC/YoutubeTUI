package views

import (
	"ytt/components"
	"ytt/daemon"
	"ytt/themes"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

type PlaylistModel struct {
	list          components.List
	width, height int
}

func Playlist() PlaylistModel {
	var rows []components.ListEntry
	for _, p := range daemon.Playlists {
		var r components.ListEntry
		r.Name = p.Title
		r.Desc = p.Channel
		r.CustomData = p.ID
		rows = append(rows, r)
	}
	list := components.NewList(rows[:], "Playlists")
	return PlaylistModel{list: list}
}

func (m PlaylistModel) Update(msg tea.Msg) (PlaylistModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		if e, ok := m.list.Hovered(); ok {
			if msg.String() == "enter" {
				ReloadTracks(e.CustomData.(string))
				cmd = Goto(ViewTracks)
				return m, cmd
			}
		}
	case tea.MouseMsg:
		if e, ok := m.list.MouseHovered(msg); ok {
			if msg.Mouse().Button == tea.MouseLeft {
				ReloadTracks(e.CustomData.(string))
				cmd = Goto(ViewTracks)
				return m, cmd
			}
		}
	}
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}
func (m PlaylistModel) View() string {
	var o string
	t := themes.Active()
	listStyle := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		PaddingLeft(2).
		Background(t.Background)
	o += listStyle.Render(m.list.View())
	return o
}
