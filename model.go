package main

// A simple program that opens the alternate screen buffer and displays mouse
// coordinates and events.

import (
	"fmt"
	"ytt/components"

	tea "github.com/charmbracelet/bubbletea"
)

func Model() tea.Model {
	rows := []string{}
	for i := 0; i < 100; i++ {
		rows = append(rows, fmt.Sprintf("Row %03d | Column A | Column B", i))
	}
	return model{
		table: components.NewTable(rows),
	}
}

type model struct {
	table components.Table
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	
	return m, cmd
}

func (m model) View() string {
	s := "Simple Table (q to quit)\n\n"
	s += m.table.View()
	return s
}
