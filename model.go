package main

// A simple program that opens the alternate screen buffer and displays mouse
// coordinates and events.

import (
	"os"
	"ytt/components"
	"ytt/daemon"
	"ytt/themes"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"
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
}

func (m model) Init() tea.Cmd {
	ok := themes.Activate("GitHub Dark")
	if !ok {
		os.Exit(1)
	}
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			themes.ActiveID = (themes.ActiveID + 1) % len(themes.Themes)
		case "q", "ctrl+c":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	}
	m.table, cmd = m.table.Update(msg)

	return m, cmd
}
func (m model) View() string {
	t := themes.Active()

	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Background(t.Background).
		Render(m.table.View())
}
