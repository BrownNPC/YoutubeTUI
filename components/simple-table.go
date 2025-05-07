package components

import (
	"ytt/themes"

	tea "github.com/charmbracelet/bubbletea/v2"
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
				if m.Cursor >= -m.Offset+m.ViewHeight -1{
					m.Offset++
				}
			}
		}

	}
	return m, nil
}
func (m List) View() string {
	var o string
	t := themes.Active()
	end := min(m.Offset+m.ViewHeight, len(m.Data))
	var visible = m.Data[m.Offset : end-1]
	base := lipgloss.NewStyle()

	// Title
	o += base.Foreground(t.BrightPurple).
		Background(t.Background).
		MarginTop(2).MarginBottom(1).
		Render(m.Title) + "\n"
	for i, e := range visible {
		i +=m.Offset
		// Name
		if i == m.Cursor{
			o+="this\n"
		}
		o += base.Foreground(t.Foreground).
			Background(t.Background).
			Render(e.Name) + "\n"
		// Description
		o += " " + base.Foreground(t.BrightGreen).
			Background(t.Background).
			MarginBottom(1).
			Render(e.Desc) + "\n"

	}

	return base.PaddingLeft(2).
		Render(o)
}
