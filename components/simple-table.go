package components

import (
	"ytt/themes"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TableEntry represents a single row in the table.
// Title is the main label, and Description provides additional context.
type TableEntry struct {
	Title       string
	Description string
}

// Table is a scrollable list model for Bubble Tea.
// Data holds all rows, Cursor is the selected index,
// ViewportH is visible rows count, and ScrollTop is the first visible index.
type Table struct {
	Title string
	Data  []TableEntry // all rows

	Selected   int // selected row index
	ViewportH  int // number of rows that fit in view
	TopVisible int // index of topmost visible row
}

// NewTable initializes a Table with entries and default viewport height.
func NewTable(data []TableEntry, title string) Table {
	return Table{Data: data, ViewportH: 10, Title: title}
}

// Update handles input and window resize messages.
// Uses a pointer receiver so mutations persist.
func (m Table) Update(msg tea.Msg) (Table, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// switch msg.String() {
		// case "up", "k":
		// 	if m.Selected >= 1 {
		// 		m.Selected--
		// 	}
		// case "down", "j":
		// 	if m.Selected < len(m.Data) {
		// 		m.Selected++
		// 	}
		// case "q", "ctrl+c":
		// 	return m, tea.Quit
		// }
		// Adjust scroll to keep cursor visible
		// if m.Selected < m.TopVisible {
		// 	m.TopVisible = m.Selected
		// }
		// if m.Selected >= m.TopVisible+m.ViewportH {
		// 	m.TopVisible = m.Selected - m.ViewportH + 2
		// }

	case tea.WindowSizeMsg:
		// Reserve one line for status/footer
		m.ViewportH = (msg.Height * 80) / 100
	}
	return m, nil
}

// View renders the visible rows as a formatted string.
func (m Table) View() string {
	var o string
	t := themes.Active()
	var style = lipgloss.NewStyle()

	title := style.Foreground(t.BrightRed).Render(m.Title)
	bottom := min(m.TopVisible+m.ViewportH, len(m.Data)-1)

	visibleRows := m.Data[m.TopVisible:bottom]
	for i, row := range visibleRows {
		cursor := "  "
		if i == m.Selected {
			cursor = "➤ "
		}
		o += cursor

		o += style.Foreground(t.BrightCyan).Bold(true).Render(row.Title) + "\n"
	}
	return lipgloss.JoinVertical(0.0, title, o)

}
