package components

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TableEntry represents a single row in the table.
// Title is the main label, and Description provides additional context.
var docStyle = lipgloss.NewStyle()

type TableEntry struct {
	Name string
	Desc string
}

func (i TableEntry) Title() string       { return i.Name }
func (i TableEntry) Description() string { return i.Desc }
func (i TableEntry) FilterValue() string { return i.Name + " " + i.Desc }

// // List is a scrollable list model for Bubble Tea.
// // Data holds all rows, Cursor is the selected index,
// // ViewportH is visible rows count, and ScrollTop is the first visible index.
type List struct {
	Title string
	list  list.Model // use table under the hood
}

var breakline = 'ㅤ' //hangul filler used to bypass not having newlines
// NewTable initializes a Table with entries and default viewport height.
func NewList(data []TableEntry, title string) List {
	var items = make([]list.Item, len(data))
	for i := range data {
		items[i] = data[i]
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	return List{Title: title, list: l}
}

func (m List) Update(msg tea.Msg) (List, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	}
	m.list, _ = m.list.Update(msg)
	return m, nil
}
func (m List) View() string {
	var o string
	o += m.list.View()
	return o
}
