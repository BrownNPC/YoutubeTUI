package views

import (
	"fmt"
	daemon "ytt/YoutubeDaemon"
	"ytt/components"
	"ytt/helpers"
	"ytt/themes"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

type PlaylistModel struct {
	list           components.List
	width, height  int
	showingOptions bool
}

var Menu = struct {
	Options        []string
	selectedOption int
	prefix         string
}{
	Options: []string{
		"Play",
		"View tracks",
	},
	prefix: "playlistMenu",
}

func RenderOptions() string {
	t := themes.Active()
	var o string
	for i, opt := range Menu.Options {
		if Menu.selectedOption == i {
			opt = lipgloss.NewStyle().
				Background(t.Background).
				Foreground(t.CursorColor).
				Faint(true).
				Render(opt)
		} else {
			opt = lipgloss.NewStyle().
				Background(t.Background).
				Foreground(t.Foreground).
				Render(opt)
		}
		opt = zone.Mark(fmt.Sprint(Menu.prefix, i), opt)
		if i != len(Menu.Options)-1 {
			opt += "\n"
		}
		o += opt
	}
	base := lipgloss.NewStyle()
	return base.
		Border(lipgloss.RoundedBorder()).
		BorderBackground(t.Background).
		BorderForeground(t.SelectionBackground).
		PaddingLeft(1).
		PaddingRight(1).
		AlignHorizontal(lipgloss.Center).
		Background(t.Background).
		Render(o)
}
func Playlist() PlaylistModel {
	var rows []components.ListEntry
	for _, p := range daemon.RegisteredPlaylists() {
		var r components.ListEntry
		r.Name = p.Title
		r.Desc = p.Channel
		r.CustomData = p.ID
		rows = append(rows, r)
	}
	list := components.NewList(rows[:], "Playlists")
	return PlaylistModel{list: list}
}
func updateMenu(keyCode rune) {
	switch keyCode {
	case tea.KeyDown, 'j':
		Menu.selectedOption++
	case tea.KeyUp, 'k':
		Menu.selectedOption--
	}
	Menu.selectedOption %= len(Menu.Options)
	if Menu.selectedOption<0{
		Menu.selectedOption=len(Menu.Options)-1
	}
}
func (m PlaylistModel) Update(msg tea.Msg) (PlaylistModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.Key().Code {
		case tea.KeyEsc:
			m.showingOptions = false
			return m, nil
		case tea.KeyEnter:
			m.showingOptions = true
			return m, nil
		default:
			if m.showingOptions {
				updateMenu(msg.Key().Code)
			}
		}
		// if e, ok := m.list.Hovered(); ok {
		// 	if msg.String() == "enter" {
		// 		ReloadTracks(e.CustomData.(string))
		// 		cmd = Goto(ViewTracks)
		// 		return m, cmd
		// 	}
		// }
	case tea.MouseMsg:
		// if e, ok := m.list.MouseHovered(msg); ok {
		// 	if msg.Mouse().Button == tea.MouseLeft {
		// 		ReloadTracks(e.CustomData.(string))
		// 		cmd = Goto(ViewTracks)
		// 		return m, cmd
		// 	}
		// }
	}
	if !m.showingOptions {
		m.list, cmd = m.list.Update(msg)
	}
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
	if m.showingOptions {
		o, _ = helpers.OverlayCenter(o, RenderOptions(), true)
	}
	return o
}
