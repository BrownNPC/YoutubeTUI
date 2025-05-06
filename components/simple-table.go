package components

import (
	"fmt"
	"github.com/charmbracelet/bubbletea"
)

// TableEntry represents a single row in the table.
// Title is the main label, and Description provides additional context.
type TableEntry struct {
	Title       string
	Description string
}

// Table is a simple scrollable table model for Bubble Tea.
// Data holds the rows, Cursor tracks the selected row,
// ViewportH is the number of visible rows, and ScrollTop
// is the index of the first visible row.
type Table struct {
	Data      []TableEntry // all rows in the table
	Cursor    int          // index of the currently selected row
	ViewportH int          // height of the viewport in rows
	ScrollTop int          // index of the topmost visible row
}

// NewTable creates a Table with the provided entries.
func NewTable(data []TableEntry) Table {
	return Table{Data: data}
}

// Update handles keyboard and window size messages.
// It moves the cursor, adjusts scrolling, and quits on q/Ctrl+C.
func (m Table) Update(msg tea.Msg) (Table, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			// Move cursor up
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down", "j":
			// Move cursor down
			if m.Cursor < len(m.Data)-1 {
				m.Cursor++
			}
		case "q", "ctrl+c":
			// Quit program
			return m, tea.Quit
		}
		// After moving the cursor, adjust scroll to keep the cursor in view
		if m.Cursor < m.ScrollTop {
			m.ScrollTop = m.Cursor
		} else if m.Cursor >= m.ScrollTop+m.ViewportH {
			m.ScrollTop = m.Cursor - m.ViewportH + 1
		}

	case tea.WindowSizeMsg:
		// Recalculate viewport height when the window size changes
		// Reserve one line for any footer/status
		m.ViewportH = msg.Height - 1

	}
	return m, nil
}

// View renders the visible portion of the table as a string.
func (m Table) View() string {
	var s string
	// Calculate end index without going past data length
	start := m.ScrollTop
	end := min(start+m.ViewportH, len(m.Data))

	// Render each visible row, prefixing the cursor
	for i := start; i < end; i++ {
		cursor := "  "
		if i == m.Cursor {
			cursor = "➤ " // highlight selected row
		}
		row := m.Data[i]
		// Combine title and description with padding
		s += fmt.Sprintf("%s%s - %s\n", cursor, row.Title, row.Description)
	}
	return s
}
