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

// does not play anything, it just allows you to add tracks to queue
type TracksModel struct {
	initialized   bool
	width, height int
	showingMenu   bool
	list          components.List
}
type ReinitTracksModelMsg struct {
	Playlist daemon.Playlist
}

func NewTracksModel(p daemon.Playlist) TracksModel {
	title := p.Title
	var rows []components.ListEntry
	for _, t := range p.Tracks {
		var r components.ListEntry
		r.Name = t.Title
		r.Desc = t.Uploader
		r.CustomData = t
		rows = append(rows, r)
	}
	list := components.NewList(rows, title)
	return TracksModel{list: list}
}

// left click menu
var TracksMenu = struct {
	Options        []string
	selectedOption int
	prefix         string
	selectedTrack  *daemon.Track // which track is this menu for?
	openedAt       image.Point   // coordinates of where we should open the menu. zero value = open in center
}{
	Options: []string{
		"Play",
		"Add to queue",
	},
	prefix: "tracksMenu",
}

func updateTracksMenuByReadingKeyboard(keyCode rune) {
	switch keyCode {
	case tea.KeyDown, 'j':
		TracksMenu.selectedOption++
	case tea.KeyUp, 'k':
		TracksMenu.selectedOption--
	}
	TracksMenu.selectedOption %= len(TracksMenu.Options)
	if TracksMenu.selectedOption < 0 {
		TracksMenu.selectedOption = len(TracksMenu.Options) - 1
	}
}
func (m TracksModel) Update(msg tea.Msg) (TracksModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.Key().Code {
		case tea.KeyEsc:
			m.showingMenu = false
			TracksMenu.selectedTrack = &daemon.Track{}
			TracksMenu.openedAt = image.Point{}
			return m, nil
		case tea.KeyEnter:
			if !m.showingMenu {
				t, ok := m.list.Hovered()
				if ok {
					m.showingMenu = true
					TracksMenu.selectedTrack = t.CustomData.(*daemon.Track)
				}
			} else {
				opt := TracksMenu.Options[TracksMenu.selectedOption]
				if opt == "Play" {
					go daemon.PlayTrack(TracksMenu.selectedTrack)
					m.showingMenu = false
				}
			}
			return m, nil
		default:
			if m.showingMenu {
				updateTracksMenuByReadingKeyboard(msg.Key().Code)
			}
		}
	case tea.MouseClickMsg:
		if m.showingMenu || msg.Mouse().Button == tea.MouseLeft || msg.Mouse().Button == tea.MouseRight {
			z := zone.Get("playlistModal")
			// hide if clicking outside modal
			if !helpers.ZoneCollision(z, msg) {
				m.showingMenu = false
				TracksMenu.selectedTrack = &daemon.Track{}
				TracksMenu.openedAt = image.Point{}
			}
		}
		// open the modal for the clicked playlist
		if e, ok := m.list.MouseHovered(msg); ok {
			hoveredTrack := e.CustomData.(*daemon.Track)
			p, _ := m.list.Hovered()
			if p.Name == e.Name && p.Desc == e.Desc {
				m.showingMenu = true
				TracksMenu.selectedTrack = hoveredTrack
				TracksMenu.openedAt.X, TracksMenu.openedAt.Y = msg.X, msg.Y
			}
		}
	case tea.MouseMsg:
		for i := range TracksMenu.Options {
			z := zone.Get(fmt.Sprint(TracksMenu.prefix, i))
			if helpers.ZoneCollision(z, msg) {
				TracksMenu.selectedOption = i
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
func (m TracksModel) View() string {
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
		if TracksMenu.openedAt.Eq(image.Point{}) {
			o, _ = helpers.OverlayCenter(o, RenderTracksMenuOptions(), true)
		} else { // draw at coordinates
			x, y := TracksMenu.openedAt.X, TracksMenu.openedAt.Y
			o = helpers.PlaceOverlay(x, y, RenderTracksMenuOptions(), o)
		}
	}
	return o
}
func RenderTracksMenuOptions() string {
	t := themes.Active()
	var o string
	for i, opt := range TracksMenu.Options {
		if TracksMenu.selectedOption == i {
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
		opt = zone.Mark(fmt.Sprint(TracksMenu.prefix, i), opt)
		if i != len(TracksMenu.Options)-1 {
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
