package components

import (
	"ytt/themes"

	"github.com/charmbracelet/lipgloss"
)

type TableEntry struct {
	Name string
	Desc string
}

// // List is a scrollable list model for Bubble Tea.
// // Data holds all rows, Cursor is the selected index,
// // ViewportH is visible rows count, and ScrollTop is the first visible index.
type List struct {
	width      int
	Title      string
	Data       []TableEntry
	ViewHeight int
	Offset     int
	Cursor     int
}

// NewTable initializes a Table with entries and default viewport height.
func NewList(data []TableEntry, title string) List {
	return List{Data: data, Title: title, ViewHeight: 20}
}

func (m List) Update(msg tea.Msg) (List, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.Cursor > 0 {
				m.Cursor--
				if m.Cursor < m.Offset {
					m.Offset--
				}
			}
		case "down":
			if m.Cursor < len(m.Data)-1 {
				m.Cursor++
				if m.Cursor >= -m.Offset+m.ViewHeight-1 {
					m.Offset++
				}
			}
		}

	}
	return m, nil
}
func (m List) View() string {
	// fields are lipgloss.Color
	var t themes.Theme = themes.Active()

	base := lipgloss.NewStyle().
		Background(t.Background).
		Width(m.width) // set to be screen width

	// Title
	title := base.
		Foreground(t.BrightPurple).
		MarginBottom(1).
		Render(m.Title)
	var listContent string
	for i, e := range m.Data {
		i += m.Offset

		content := lipgloss.JoinVertical(
			lipgloss.Bottom,
			// Name
			base.
				Background(t.Blue).
				PaddingLeft(1).
				Foreground(t.BrightGreen).
				Background(t.BrightWhite).
				Render(e.Name),
			// Description
			base.
				MarginBottom(1).
				PaddingLeft(2).
				Render(e.Desc),
		)
		listContent += content

	}
	content := lipgloss.JoinVertical(
		lipgloss.Top,
		title, listContent,
	)
	return base.
		Render(content)
}
