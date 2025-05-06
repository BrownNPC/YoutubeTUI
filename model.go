package main

// A simple program that opens the alternate screen buffer and displays mouse
// coordinates and events.

import (
	"fmt"
	"ytt/components"
	"ytt/themes"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func Model() tea.Model {
	var rows [100]components.TableEntry
	for i, t := range rows {
		t.Title = fmt.Sprintf("Row %d", i)
		t.Description = fmt.Sprintf("Description of row %d", i)
		rows[i] = t
	}
	return model{
		table: components.NewTable(rows[:], "Playlists"),
	}
}

type model struct {
	table components.Table
}

func (m model) Init() tea.Cmd {
	themes.Activate(1)
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "q", "ctrl+c":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
	}
	m.table, cmd = m.table.Update(msg)

	return m, cmd
}

func (m model) View() string {
	t := themes.Active()
	var style lipgloss.Style
	style.Background(t.Background)
	s := "Simple Table (q to quit)\n\n"

	return lipgloss.JoinVertical(0.2,
		style.Render(s),
		m.table.View(),
	)
}
