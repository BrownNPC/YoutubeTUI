package main

// A simple program that opens the alternate screen buffer and displays mouse
// coordinates and events.

import (
	"time"
	menu "ytt/globalmenu"
	"ytt/helpers"
	"ytt/views"

	tea "github.com/charmbracelet/bubbletea/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

func Model() tea.Model {

	return model{
		view:            views.ViewPlaylists,
		playlistView:    views.Playlist(),
		changeThemeView: views.ChangeTheme(),
		menuOpened:      true,
		openAtCenter:    true,
	}
}

type model struct {
	playlistView    views.PlaylistModel
	changeThemeView views.ChangeThemeModel

	width, height    int
	view             views.ViewMsg
	menuOpened       bool
	openAtCenter     bool
	openatX, openatY int
}

func (m model) Init() tea.Cmd {
	return tea.Every(time.Millisecond*16, func(t time.Time) tea.Msg {
		return TickMsg{}
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		menu.Update(msg)
		m.playlistView, cmd = m.playlistView.Update(msg)
		cmd = views.ActiveTracksModel.Update(msg)
		m.changeThemeView, cmd = m.changeThemeView.Update(msg)

	case TickMsg:
		return m, CmdTick
	case tea.MouseClickMsg:
		if msg.Button == tea.MouseRight {
			m.openAtCenter = false
			m.menuOpened = !m.menuOpened
			m.openatX = msg.X
			m.openatY = msg.Y
		} else if msg.Button == tea.MouseLeft && m.menuOpened {
			z := zone.Get("menu")
			if !helpers.ZoneCollision(z, msg) {
				m.menuOpened = false
			}
		}
		// Playlists view has to close the playlist click menu [views.PlaylistMenu]
		m.updateViews(msg)
	case tea.KeyMsg:
		switch msg.String() {
		case " ", "space":
			m.openAtCenter = true
			m.menuOpened = !m.menuOpened
		case "esc":
			m.menuOpened = false
		case "q", "ctrl+c":
			return m, tea.Quit
		}

	case views.ViewMsg:
		m.view = msg
		m.menuOpened = false
		return m, cmd
	}
	if !m.menuOpened {
		cmd = m.updateViews(msg)
	} else {
		cmd = menu.Update(msg)
	}

	return m, cmd
}
func (m *model) updateViews(msg tea.Msg) (cmd tea.Cmd) {
	switch m.view {
	case views.ViewPlaylists:
		m.playlistView, cmd = m.playlistView.Update(msg)
	case views.ViewTracks:
		cmd = views.ActiveTracksModel.Update(msg)
	case views.ViewChangeTheme:
		m.changeThemeView, cmd = m.changeThemeView.Update(msg)
	}
	return
}
func (m model) visibleView() string {
	var content string

	switch m.view {
	case views.ViewPlaylists:
		content = m.playlistView.View()
	case views.ViewChangeTheme:
		content = m.changeThemeView.View()
	case views.ViewTracks:
		views.ActiveTracksModel.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
		content = views.ActiveTracksModel.View()

	}
	return content
}
func (m model) View() (view string) {
	content := m.visibleView()
	view = content
	// view, _ = helpers.Overlay(view, content, 0, 0, true)
	if m.menuOpened { // render menu as an overlay
		if m.openAtCenter {
			z := zone.Get("activeList")

			msg := tea.MouseMotionMsg{X: m.openatX, Y: m.openatY}
			if helpers.ZoneCollision(z, msg) {
				m.openatX = z.EndX
				m.openatY = z.EndY
			}
			view, _ = helpers.OverlayCenter(view, menu.View(false), false)
		} else {
			view = helpers.PlaceOverlay(m.openatX, m.openatY, menu.View(true), view)
		}
	}
	return zone.Scan(view)
}
