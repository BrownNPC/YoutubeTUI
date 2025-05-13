package components

import (
	"ytt/themes"

	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TableEntry struct {
	Name string
	Desc string
}

// // List is a paginated list for Bubble Tea.
type List struct {
	paginator  paginator.Model
	width      int
	Title      string
	Data       []TableEntry
	ViewHeight int
	Cursor     int // relative to visible rows
}

// NewTable initializes a Table with entries and default viewport height.
func NewList(data []TableEntry, title string) List {
	var pag = paginator.New()
	pag.Type = paginator.Dots
	pag.ActiveDot = "◉"
	pag.SetTotalPages(len(data))
	return List{Data: data, Title: title, paginator: pag}
}

func (m List) Update(msg tea.Msg) (List, tea.Cmd) {
	m.paginator.SetTotalPages(len(m.Data))
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		msg.Height = max(msg.Height, 1)
		m.ViewHeight = max((msg.Height*80)/100, 1)
		m.paginator.PerPage = max((m.ViewHeight*35)/100, 1)
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			m.paginator.PrevPage()
		case "right", "l":
			m.paginator.NextPage()
		case "up", "k":
			m.Cursor--
		case "down", "j":
			if m.Cursor < len(m.Data)-1 {
				m.Cursor++
			}
		}
	}
	// Go to next page if end of current page (down)
	if m.Cursor >= m.paginator.ItemsOnPage(len(m.Data)) {
		// go to next page
		if m.paginator.OnLastPage() {
			m.Cursor = m.paginator.ItemsOnPage(len(m.Data)) - 1
		} else {
			m.paginator.NextPage()
			m.Cursor = 0
		}

	}
	if m.Cursor < 0 {
		if m.paginator.Page > 0 {
			m.paginator.PrevPage()
			m.Cursor = m.paginator.ItemsOnPage(len(m.Data)) - 1
		} else {
			m.Cursor = 0
		}
	}
	m.Cursor %= len(m.Data) + 1
	return m, cmd
}

func (m List) View() string {
	// fields are lipgloss.Color
	var t themes.Theme = themes.Active()
	var accentColor, selectionColor = themes.AccentColor(), themes.SelectionColor()

	base := lipgloss.NewStyle().
		Background(t.Background)
	// Title
	title := base.
		Foreground(t.Red).
		Render(m.Title)
	var listContent string
	start, end := m.paginator.GetSliceBounds(len(m.Data))
	start = min(start, end-1)
	for i, e := range m.Data[start:end] {
		var selected string = " "
		var nameColor, descColor = accentColor, t.Foreground

		if i == m.Cursor {
			selected = lipgloss.NewStyle().Foreground(t.CursorColor).Render("│")
			nameColor = selectionColor
		}
		if len(e.Name) > 40 {
			e.Name = e.Name[:40] + "…"
		}
		listContent += selected + base.
			Foreground(nameColor).
			Render(e.Name) + "\n"
		// Description
		listContent += selected + base.
			MarginBottom(1).
			Foreground(descColor).
			Render(e.Desc) + "\n"
	}
	return base.
		PaddingTop(1).
		Render(title + "\n" + m.paginator.View() + "\n\n" + listContent)
}
