package views

import (
	"fmt"
	"ytt/components"
	"ytt/daemon"

	tea "github.com/charmbracelet/bubbletea/v2"
)

func ReloadTracks(id string) {
	m := ActiveTracksModel
	playlist, ok := daemon.Playlists[id]
	if ok {
		m.initialized = true
		var tracks = make([]components.ListEntry, 0, len(playlist.Entries))

		for _, t := range playlist.Entries {
			var e components.ListEntry
			e.Name = t.Title
			e.Desc = fmt.Sprintf("%s·%s", t.Uploader, t.DurationString)
			e.CustomData = t.ID
			tracks = append(tracks, e)
		}
		var l = components.NewList(tracks, playlist.Title)
		m.list = l
	}
	ActiveTracksModel = m
}

type TracksModel struct {
	initialized        bool
	width, height      int
	loadedTracksFromId string

	list components.List
}

var ActiveTracksModel = TracksModel{}
func init(){
	ReloadTracks("")
}
func (m *TracksModel) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.list, _ = m.list.Update(msg)
	case tea.KeyMsg:
		// if e, ok := m.list.Hovered(); ok {
		// 	if msg.String() == "enter" {
		// 		SelectedTrackId = e.CustomData.(string)
		// 		cmd = Goto(ViewTracks)
		// 		return m, cmd
		// 	}
		// }
	case tea.MouseMsg:
		// if e, ok := m.list.MouseHovered(msg); ok {
		// 	if msg.Mouse().Button == tea.MouseLeft {
		// 		SelectedTrackId = e.CustomData.(string)
		// 		cmd = Goto(ViewTracks)
		// 		return m, cmd
		// 	}
		// }
	}
	if m.initialized {
		m.list, cmd = m.list.Update(msg)
	}
	return cmd
}
func (m TracksModel) View() string {
	var o string
	if !m.initialized {
		return "Please select a playlist first"
	}

	o += m.list.View()
	return o
}
