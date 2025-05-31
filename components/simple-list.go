package components

import (
	"regexp"
	"strings"
	"ytt/themes"

	"github.com/charmbracelet/bubbles/v2/paginator"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

// ListEntry represents a single item in the list with name and description
type ListEntry struct {
	Name       string
	Desc       string
	CustomData any
}

// List implements a paginated, searchable list component
type List struct {
	isSearching   bool            // Whether search input is focused
	searchChanged bool            // Flag for search filter updates
	SearchQuery   string          // Current search query text
	paginator     paginator.Model // Handles pagination state
	width, height int             // Component width
	Title         string          // Header title for the list
	AllData       []ListEntry     // Complete unfiltered dataset
	FilteredData  []ListEntry     // Data filtered by search query
	start, end    int             // Slice bounds for current page
	ViewHeight    int             // Total available height for rendering
	Cursor        int             // Current selection position (relative to visible page)
	SelectedName  string          // name of the selected element, set by consumer, selected elements will blink

}

var alphaNum = regexp.MustCompile(`^[a-zA-Z0-9]$`)

// NewList creates a new list component with given data and title
func NewList(data []ListEntry, title string) List {
	// Configure pagination defaults
	var pag = paginator.New()
	pag.Type = paginator.Dots
	pag.ActiveDot = "◉"
	pag.SetTotalPages(len(data))

	// Configure search input
	return List{AllData: data, Title: title, paginator: pag}
}

// Update handles messages and updates component state
func (m List) Update(msg tea.Msg) (List, tea.Cmd) {
	var cmd tea.Cmd
	// Update filtered data when not searching or when search query changes
	if !m.isSearching {
		// Show full dataset when not searching
		m.paginator.SetTotalPages(len(m.AllData))
		m.FilteredData = m.AllData
	} else if m.searchChanged {
		m.searchChanged = false
		// Filter data based on search query
		var selectedData []ListEntry
		for _, e := range m.AllData {
			if m.SearchQuery != "" && !strings.Contains(strings.ToLower(e.Name+e.Desc), strings.ToLower(m.SearchQuery)) {
				continue
			}
			selectedData = append(selectedData, e)
		}
		m.paginator.SetTotalPages(len(selectedData))
		m.FilteredData = selectedData
	}

	// Handle different message types
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// Adjust dimensions based on window size
		m.width = msg.Width
		m.height = msg.Height
		msg.Height = max(msg.Height, 1)
		m.ViewHeight = max((msg.Height*80)/100, 1)
		m.paginator.PerPage = max((m.ViewHeight*35)/100, 1)
	case tea.MouseWheelMsg:
		switch msg.Button {
		case tea.MouseWheelUp:
			m.paginator.PrevPage()
		case tea.MouseWheelDown:
			m.paginator.NextPage()
		}
	case tea.MouseMsg:

	case tea.KeyMsg:
		switch msg.String() {
		case "/": // Start search
			m.isSearching = !m.isSearching
			if m.isSearching {
				m.SearchQuery = ""
				return m, cmd
			}
		case "esc": // Finish search
			if m.isSearching {
				m.isSearching = false
			}
		case "left", "h": // Previous page
			m.paginator.PrevPage()
		case "right", "l": // Next page
			m.paginator.NextPage()
		case "up", "k": // Move selection up
			m.Cursor--
		case "down", "j": // Move selection down
			if m.Cursor < len(m.AllData)-1 {
				m.Cursor++
			}
		}
	}

	// Handle cursor position wrapping between pages
	if m.Cursor >= m.paginator.ItemsOnPage(len(m.FilteredData)) {
		// Move to next page if at end of current page
		if m.paginator.OnLastPage() {
			m.Cursor = m.paginator.ItemsOnPage(len(m.FilteredData)) - 1
		} else {
			m.paginator.NextPage()
			m.Cursor = 0
		}
	}
	if m.Cursor < 0 {
		// Move to previous page if at start of current page
		if m.paginator.Page > 0 {
			m.paginator.PrevPage()
			m.Cursor = m.paginator.ItemsOnPage(len(m.FilteredData)) - 1
		} else {
			m.Cursor = 0
		}
	}
	m.Cursor %= len(m.FilteredData) + 1 // Ensure cursor stays in valid range

	// Update search input if active
	if msg, ok := msg.(tea.KeyMsg); ok && m.isSearching {
		if len(msg.String()) == 1 && alphaNum.MatchString(msg.String()) {
			m.SearchQuery += msg.String()
			m.searchChanged = true
		}
		if msg.Key().Code == tea.KeyBackspace {
			runes := []rune(m.SearchQuery)
			if len(runes) > 0 {
				m.SearchQuery = string(runes[:len(runes)-1])
				m.searchChanged = true
			}
		}
	}
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
		zoneId := e.Name // name is truncated later
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
		element += selected + base.
			Foreground(nameColor).
			Blink(m.SelectedName == displayName).
			Render(displayName)
		element = zone.Mark(zoneId, element) + "\n"
		if e.Desc != "" {
			element += selected + base.
				MarginBottom(1).
				Foreground(descColor).
				Render(e.Desc) + "\n"
		} else {
			element += "\n"
		}
		listContent += element
	}

	// Show search input if active
	var search string
	if m.isSearching {
		// m.input.Styles.Focused.Text = m.input.Styles.Focused.Text.
		// 	Foreground(t.Foreground).Background(t.Background)
		search = "> " + m.SearchQuery + "\n"
	}
	// Combine all components
	view := base.
		PaddingTop(1).
		Border(HangulFillerBorder()).
		BorderBackground(t.Background).
		BorderForeground(t.Foreground).
		Render(title + "\n" + search + m.paginator.View() + "\n\n" + listContent)
	return zone.Mark("activeList", view)
}
func (m *List) MouseHovered(msg tea.MouseMsg) (ListEntry, bool) {
	start, end := m.paginator.GetSliceBounds(len(m.FilteredData))
	mouse := msg.Mouse()
	if len(m.FilteredData) == 0 {
		return ListEntry{}, false
	}
	for i, e := range m.FilteredData[start:end] {
		z := zone.Get(e.Name)
		if z.IsZero() {
			break
		}
		if mouse.X >= z.StartX && mouse.X <= z.EndX &&
			mouse.Y >= z.StartY && mouse.Y <= z.EndY {
			m.Cursor = i
			return e, true
		}
	}
	return ListEntry{}, false
}
func (m List) Hovered() (ListEntry, bool) {
	start, end := m.paginator.GetSliceBounds(len(m.FilteredData))
	if len(m.FilteredData) == 0 {
		return ListEntry{}, false
	}
	return m.FilteredData[start:end][m.Cursor], true
}
