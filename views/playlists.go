package views

import (
	"fmt"
	"image"
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
	showingMenu bool
}

// left click menu
var PlaylistMenu = struct {
	Options          []string
	selectedOption   int
	prefix           string
	selectedPlaylist daemon.Playlist // which playlist is this menu for?
	openedAt         image.Point     // coordinates of where we should open the menu. zero value = open in center
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
	for i, opt := range PlaylistMenu.Options {
		if PlaylistMenu.selectedOption == i {
			opt = lipgloss.NewStyle().
				Background(t.Background).
				Foreground(t.CursorColor).
				Bold(true).
				Render(opt)
		} else {
			opt = lipgloss.NewStyle().
				Background(t.Background).
				Foreground(t.Foreground).
				Faint(true).
				Render(opt)
		}
		opt = zone.Mark(fmt.Sprint(PlaylistMenu.prefix, i), opt)
		if i != len(PlaylistMenu.Options)-1 {
			opt += "\n"
		}
		o += opt
	}
	base := lipgloss.NewStyle()
	o = base.
		Border(lipgloss.RoundedBorder()).
		BorderBackground(t.Background).
		BorderForeground(t.SelectionBackground).
		PaddingLeft(1).
		PaddingRight(1).
		AlignHorizontal(lipgloss.Center).
		Background(t.Background).
		Render(o)
	return zone.Mark("playlistModal", o)
}
func Playlist() PlaylistModel {
	var rows []components.ListEntry
	for _, p := range daemon.RegisteredPlaylists() {
		var r components.ListEntry
		r.Name = p.Title
		r.Desc = p.Channel
		r.CustomData = p
		rows = append(rows, r)
	}
	list := components.NewList(rows[:], "Playlists")
	return PlaylistModel{list: list}
}
func updateMenuByReadingKeyboard(keyCode rune) {
	switch keyCode {
	case tea.KeyDown, 'j':
		PlaylistMenu.selectedOption++
	case tea.KeyUp, 'k':
		PlaylistMenu.selectedOption--
	}
	PlaylistMenu.selectedOption %= len(PlaylistMenu.Options)
	if PlaylistMenu.selectedOption < 0 {
		PlaylistMenu.selectedOption = len(PlaylistMenu.Options) - 1
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
			m.showingMenu = false
			PlaylistMenu.selectedPlaylist = daemon.Playlist{}
			PlaylistMenu.openedAt = image.Point{}
			return m, nil
		case tea.KeyEnter:
			if !m.showingMenu {
				p, ok := m.list.Hovered()
				if ok {
					m.showingMenu = true
					PlaylistMenu.selectedPlaylist = p.CustomData.(daemon.Playlist)
				}
			}
			return m, nil
		default:
			if m.showingMenu {
				updateMenuByReadingKeyboard(msg.Key().Code)
			}
		}
	case tea.MouseClickMsg:
		if m.showingMenu || msg.Mouse().Button == tea.MouseLeft || msg.Mouse().Button == tea.MouseRight {
			z := zone.Get("playlistModal")
			// hide if clicking outside modal
			if !helpers.ZoneCollision(z, msg) {
				m.showingMenu = false
				PlaylistMenu.selectedPlaylist = daemon.Playlist{}
				PlaylistMenu.openedAt = image.Point{}
			}
		}
		// open the modal for the clicked playlist
		if e, ok := m.list.MouseHovered(msg); ok {
			hoveredPlaylist := e.CustomData.(daemon.Playlist)
			p, _ := m.list.Hovered()
			if p.Name == e.Name && p.Desc == e.Desc {
				m.showingMenu = true
				PlaylistMenu.selectedPlaylist = hoveredPlaylist
				PlaylistMenu.openedAt.X, PlaylistMenu.openedAt.Y = msg.X, msg.Y
			}
		}
	case tea.MouseMsg:
		for i := range PlaylistMenu.Options {
			z := zone.Get(fmt.Sprint(PlaylistMenu.prefix, i))
			if helpers.ZoneCollision(z, msg) {
				PlaylistMenu.selectedOption = i
			}
		}
		// if e, ok := m.list.MouseHovered(msg); ok {
		// }
	}
	if !m.showingMenu {
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
	if m.showingMenu {
		// zero value, draw at center
		if PlaylistMenu.openedAt.Eq(image.Point{}) {
			o, _ = helpers.OverlayCenter(o, RenderOptions(), true)
		} else { // draw at coordinates
			x, y := PlaylistMenu.openedAt.X, PlaylistMenu.openedAt.Y
			o = helpers.PlaceOverlay(x, y, RenderOptions(), o)
		}
	}
	return o
}
