package main

// A simple program that opens the alternate screen buffer and displays mouse
// coordinates and events.

import (
	"os"
	"time"
	"ytt/components"
	"ytt/daemon"
	"ytt/themes"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ViewMsg int

const (
	ViewPlaylists ViewMsg = iota
	ViewSettings
)

func Model() tea.Model {
	var rows []components.TableEntry
	for _, p := range daemon.Playlists {
		for _, t := range p.Entries {
			var r components.TableEntry
			r.Name = t.Title
			r.Desc = t.Uploader
			rows = append(rows, r)

		}
	}
	return model{
		table: components.NewList(rows[:], "Playlists"),
	}
}

type model struct {
	table         components.List
	width, height int
	view          ViewMsg
}

func (m model) Init() tea.Cmd {
	ok := themes.Activate("GitHub Dark")
	if !ok {
		os.Exit(1)
	}
	return tea.Every(time.Millisecond*16, func(t time.Time) tea.Msg {
		return TickMsg{}
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case TickMsg:
		return m, CmdTick
	case tea.KeyMsg:
		switch msg.String() {
		case "1":
			m.view = ViewPlaylists
		case "0":
			m.view = ViewSettings
		case "enter":
			themes.ActiveID = (themes.ActiveID + 1) % len(themes.Themes)
		case "q", "ctrl+c":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case ViewMsg:
		switch msg {
		case ViewPlaylists:
			m.view = ViewPlaylists
		}
	}
	switch m.view {
	case ViewPlaylists:
		m.table, cmd = m.table.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}
func (m model) View() (view string) {
	t := themes.Active()
	switch m.view {
	case ViewPlaylists:
		view = lipgloss.NewStyle().
			Width(m.width).
			Height(m.height).
			PaddingLeft(2).
			Background(t.Background).
			Render(m.table.View())
	case ViewSettings:

	}
	return view
}
