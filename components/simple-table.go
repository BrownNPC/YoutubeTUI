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
	Title        string
	Data         []TableEntry
	ViewHeight   int
	CurrentPage  int
	ItemsPerPage int
	TotalPages   int
	Cursor       int // which of the visible elements is selected
}

// NewTable initializes a Table with entries and default viewport height.
func NewList(data []TableEntry, title string) List {
	return List{Data: data, Title: title, ViewHeight: 20}
}

func (m List) Update(msg tea.Msg) (List, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		}
	}
	return m, nil
}
func (m List) View() string {
	var o string
	t := themes.Active()

	pages := len(m.Data) / m.ViewHeight

	base := lipgloss.NewStyle()
	var visible = m.Data[m.Cursor:m.ViewHeight]

	// Title
	o += base.Foreground(t.BrightPurple).
		Background(t.Background).
		MarginTop(2).MarginBottom(1).
		Render(m.Title) + "\n"
	for _, e := range visible {
		// Name
		o += base.Foreground(t.Foreground).
			Render(e.Name) + "\n"
		// Description
		o += " " + base.Foreground(t.BrightGreen).
			MarginBottom(1).
			Render(e.Desc) + "\n"

	}

	return base.PaddingLeft(2).
		Render(o)
}
