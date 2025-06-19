package components

import (
	"strings"

	"ytt/helpers"
	"ytt/themes"

	"github.com/charmbracelet/bubbles/v2/paginator"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

const (
	maxNameLength = 40
)

// ListEntry represents a single item in the list
type ListEntry struct {
	Name       string
	Desc       string
	CustomData any
}

// List implements a paginated, searchable list component
type List struct {
	isSearching   bool
	searchChanged bool
	SearchQuery   string
	paginator     paginator.Model
	width, height int
	Title         string
	AllData       []ListEntry
	FilteredData  []ListEntry
	ViewHeight    int
	Cursor        int
	SelectedName  string
}

// NewList creates a new list component
func NewList(data []ListEntry, title string) List {
	p := paginator.New()
	p.Type = paginator.Dots
	p.ActiveDot = "◉"
	p.SetTotalPages(len(data))

	return List{
		AllData:   data,
		Title:     title,
		paginator: p,
	}
}

// Update handles messages and updates component state
func (m List) Update(msg tea.Msg) (List, tea.Cmd) {
	var cmd tea.Cmd
	oldPage := m.paginator.Page

	// Update paginator first
	m.paginator, cmd = m.paginator.Update(msg)
	if oldPage != m.paginator.Page {
		m.Cursor = 0
	}

	// Update filtered data
	if !m.isSearching {
		m.paginator.SetTotalPages(len(m.AllData))
		m.FilteredData = m.AllData
	} else if m.searchChanged {
		m.searchChanged = false
		m.FilteredData = m.filterData()
		m.paginator.SetTotalPages(len(m.FilteredData))
		m.paginator.Page = 0
		m.Cursor = 0
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.handleResize(msg)
	case tea.KeyMsg:
		m.handleKeyPress(msg)
	case tea.MouseWheelMsg:
		switch msg.Button {
		case tea.MouseWheelUp:
			m.paginator.PrevPage()
		case tea.MouseWheelDown:
			m.paginator.NextPage()
		}
	}
	m.adjustCursor()
	return m, cmd
}

var base = lipgloss.NewStyle()

// View renders the list component UI
func (m List) View() string {
	// Get current theme and colors
	var t themes.Theme = themes.Active()
	accentColor, selectionColor := themes.AccentColor(), themes.SelectionColor()

	// Base styling
	base = base.
		Background(t.Background)

	// Render title
	title := base.
		Foreground(themes.AccentColor()).
		Underline(true).
		Render(m.Title)
	var listContent string

	// Get current page bounds
	start, end := m.paginator.GetSliceBounds(len(m.FilteredData))
	var data = m.FilteredData
	if len(m.FilteredData) > 0 && start <= end {
		data = m.FilteredData[start:end]
	}
	// Render each item in current page
	for i, e := range data {
		var selected string = " "
		// Highlight selected item
		nameColor, descColor := accentColor, t.Foreground
		if i == m.Cursor {
			selected = lipgloss.NewStyle().Foreground(t.CursorColor).Render("│")
			nameColor = selectionColor
		}

		// Truncate long names
		displayName := e.Name
		if len(e.Name) > 40 {
			displayName = e.Name[:40] + "…"
		}
		// Render name and description
		var element string
		displayName = selected + base.
			Foreground(nameColor).
			Blink(m.SelectedName == displayName).
			Render(displayName)
		element += zone.Mark(e.Name, displayName) + "\n"
		if e.Desc != "" {
			desc := selected + base.
				MarginBottom(1).
				Foreground(descColor).
				Render(e.Desc)
			element += zone.Mark(e.Name+"desc", desc) + "\n"
		} else {
			element += "\n"
		}
		listContent += element
	}

	// Show search input if active
	var search string
	if m.isSearching {
		search = "> " + m.SearchQuery + "\n"
	}
	// Combine all components
	view := base.
		PaddingTop(1).
		Border(HangulFillerBorder()).
		BorderBackground(t.Background).
		BorderForeground(t.Foreground).
		Render(title + "\n" + search + m.paginator.View() + "\n\n" + listContent)
	return zone.Mark("activeList", view)// used by model.View
}

// check if mouse is hovering over an item
func (m *List) MouseHovered(msg tea.MouseMsg) (ListEntry, bool) {
	start, end := m.paginator.GetSliceBounds(len(m.FilteredData))
	if start >= len(m.FilteredData) {
		return ListEntry{}, false
	}

	for i, e := range m.FilteredData[start:end] {
		z, z2 := zone.Get(e.Name), zone.Get(e.Name+"desc")
		if z.InBounds(msg) || helpers.ZoneCollision(z2, msg) {
			m.Cursor = i
			return e, true
		}
	}
	return ListEntry{}, false
}

// return currently selected item
func (m List) Hovered() (ListEntry, bool) {
	start, end := m.paginator.GetSliceBounds(len(m.FilteredData))
	if start >= len(m.FilteredData) || len(m.FilteredData[start:end]) == 0 {
		return ListEntry{}, false
	}
	return m.FilteredData[start:end][m.Cursor], true
}

// apply filter to list data
func (m List) filterData() []ListEntry {
	if m.SearchQuery == "" {
		return m.AllData
	}

	var filtered []ListEntry
	query := strings.ToLower(m.SearchQuery)

	for _, e := range m.AllData {
		content := strings.ToLower(e.Name + e.Desc)
		if strings.Contains(content, query) {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

// resize events
func (m *List) handleResize(msg tea.WindowSizeMsg) {
	m.width = msg.Width
	m.height = msg.Height
	m.ViewHeight = max((msg.Height*80)/100, 1)
	m.paginator.PerPage = max((m.ViewHeight*35)/100, 1)
}

// keyboard input
func (m *List) handleKeyPress(msg tea.KeyMsg) {
	if m.isSearching && len(msg.String()) == 1 && msg.String() != "/" {
		m.SearchQuery += msg.String()
		m.searchChanged = true
		return
	}
	switch msg.String() {
	case "backspace":
		if m.isSearching && len(m.SearchQuery) > 0 {
			m.SearchQuery = m.SearchQuery[:len(m.SearchQuery)-1]
			m.searchChanged = true
		}
	case "/":
		m.isSearching = !m.isSearching
		if !m.isSearching {
			m.SearchQuery = ""
		}
	case "esc":
		m.isSearching = false
		m.SearchQuery = ""
	case "up", "k":
		m.Cursor--
	case "down", "j":
		m.Cursor++
	default:
	}
}

// adjustCursor ensures cursor stays within valid bounds
func (m *List) adjustCursor() {
	itemsPerPage := m.paginator.ItemsOnPage(len(m.FilteredData))
	if itemsPerPage == 0 {
		m.Cursor = 0
		return
	}

	switch {
	case m.Cursor < 0:
		if m.paginator.Page > 0 {
			m.paginator.PrevPage()
			m.Cursor = m.paginator.ItemsOnPage(len(m.FilteredData)) - 1
		} else {
			m.Cursor = 0
		}
	case m.Cursor >= itemsPerPage:
		if !m.paginator.OnLastPage() {
			m.paginator.NextPage()
			m.Cursor = 0
		} else {
			m.Cursor = itemsPerPage - 1
		}
	}
}

// max returns the larger of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
