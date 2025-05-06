package components

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type Table struct {
	Data      []string // your "table" rows
	Cursor    int      // current cursor position
	ViewportH int      // height of visible rows
	ScrollTop int      // where viewport starts
}

func NewTable(data []string) Table {
	return Table{Data: data}
}

func (m Table) Update(msg tea.Msg) (Table, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
				if m.Cursor < m.ScrollTop {
					m.ScrollTop--
				}
			}
		case "down", "j":
			if m.Cursor < len(m.Data)-1 {
				m.Cursor++
				if m.Cursor >= m.ScrollTop+m.ViewportH {
					m.ScrollTop++
				}
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.ViewportH = msg.Height-1
	}
	return m, nil
}
func (m Table) View() string {
	var s string

	end := min(m.ScrollTop+m.ViewportH, len(m.Data))

	for i := m.ScrollTop; i < end; i++ {
		cursor := "  "
		if i == m.Cursor {
			cursor = "➤ "
		}
		s += fmt.Sprintf("%s%s\n", cursor, m.Data[i])
	}
	return s
}
