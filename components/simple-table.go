package components

import (
	"fmt"
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

	Selected   int // index of selected element (cursor)
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
		var BottomVisible = m.ViewportH + m.TopVisible - 1
		switch msg.String() {
		case "up", "k":
			if m.Selected > 0 {
				m.Selected--
				if m.Selected == m.TopVisible && m.TopVisible-1 != -1 {
					m.TopVisible--
				}
			}
		case "down", "j":
			if m.Selected < len(m.Data)-1 {
				m.Selected++
			}
			if m.Selected == BottomVisible {
				m.TopVisible++
			}
		}
		// If selected row moves below viewport, adjust TopVisible

	case tea.WindowSizeMsg:
		m.ViewportH = (msg.Height * 80) / 100
	}
	return m, nil
}

// View renders the visible rows as a formatted string.
func (m Table) View() string {
	t := themes.Active()
	var style lipgloss.Style

	title := style.Foreground(t.BrightRed).Underline(true).Render(m.Title)
	bottom := min(m.TopVisible+m.ViewportH, len(m.Data))
	visibleRows := m.Data[m.TopVisible:bottom]

	var output string
	for i, row := range visibleRows {
		i := i + m.TopVisible // not relative
		isSelected := i == m.Selected

		// Format counter with padding and trailing space
		counter := fmt.Sprintf("%4d ", i)

		// Base styles
		counterStyle := style.Foreground(t.BrightBlack)
		lineStyle := style.Foreground(t.BrightCyan).Bold(true)

		// Apply selection background if needed
		if isSelected {
			selectedBg := t.SelectionBackground
			counterStyle = counterStyle.Background(selectedBg)
			lineStyle = lineStyle.Background(selectedBg)
		}

		output += counterStyle.Render(counter) + lineStyle.Render(row.Title) + "\n"
	}

	return lipgloss.JoinVertical(0.0, title, output)
}
