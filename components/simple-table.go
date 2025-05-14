package components

import (
	"strings"
	"ytt/themes"

	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ListEntry represents a single item in the list with name and description
type ListEntry struct {
	Name string
	Desc string
}

// List implements a paginated, searchable list component
type List struct {
	input         textinput.Model    // Search input field
	isSearching   bool               // Whether search input is focused
	searchChanged bool               // Flag for search filter updates
	SearchQuery   string             // Current search query text
	paginator     paginator.Model    // Handles pagination state
	width         int                // Component width
	Title         string             // Header title for the list
	AllData       []ListEntry       // Complete unfiltered dataset
	FilteredData  []ListEntry       // Data filtered by search query
	start, end    int                // Slice bounds for current page
	ViewHeight    int                // Total available height for rendering
	Cursor        int                // Current selection position (relative to visible page)
}

// NewList creates a new list component with given data and title
func NewList(data []ListEntry, title string) List {
	// Configure pagination defaults
	var pag = paginator.New()
	pag.Type = paginator.Dots
	pag.ActiveDot = "◉"
	pag.SetTotalPages(len(data))

	// Configure search input
	var input = textinput.New()
	input.KeyMap.CharacterBackward.Unbind()  // Disable text cursor movement
	input.KeyMap.CharacterForward.Unbind()   // Keep only line-based navigation

	return List{AllData: data, Title: title, paginator: pag, input: input}
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
		msg.Height = max(msg.Height, 1)
		m.ViewHeight = max((msg.Height*80)/100, 1)
		m.paginator.PerPage = max((m.ViewHeight*35)/100, 1)

	case tea.KeyMsg:
		switch msg.String() {
		case "/":  // Start search
			if !m.isSearching {
				m.isSearching = true
				m.input.SetValue("")
				m.SearchQuery = ""
				cmd = m.input.Focus()
				return m, cmd
			}
		case "escape":  // Finish search
			if m.isSearching {
				m.isSearching = false
			}
		case "left", "h":  // Previous page
			m.paginator.PrevPage()
		case "right", "l":  // Next page
			m.paginator.NextPage()
		case "up", "k":  // Move selection up
			m.Cursor--
		case "down", "j":  // Move selection down
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
	m.Cursor %= len(m.FilteredData) + 1  // Ensure cursor stays in valid range

	// Update search input if active
	if m.isSearching {
		m.input, cmd = m.input.Update(msg)
		if m.input.Value() != m.SearchQuery {
			m.Cursor = 0  // Reset cursor when search changes
			m.paginator.Page = 0  // Reset to first page
			m.SearchQuery = m.input.Value()
			m.searchChanged = true
		}
	}
	return m, cmd
}

// View renders the list component UI
func (m List) View() string {
	// Get current theme and colors
	var t themes.Theme = themes.Active()
	accentColor, selectionColor := themes.AccentColor(), themes.SelectionColor()

	// Base styling
	base := lipgloss.NewStyle().
		Background(t.Background)

	// Render title
	title := base.
		Foreground(t.Red).
		Render(m.Title)
	var listContent string

	// Get current page bounds
	m.start, m.end = m.paginator.GetSliceBounds(len(m.FilteredData))
	
	// Render each item in current page
	for i, e := range m.FilteredData[m.start:m.end] {
		var selected string = " "
		nameColor, descColor := accentColor, t.Foreground

		// Highlight selected item
		if i == m.Cursor {
			selected = lipgloss.NewStyle().Foreground(t.CursorColor).Render("│")
			nameColor = selectionColor
		}

		// Truncate long names
		if len(e.Name) > 40 {
			e.Name = e.Name[:40] + "…"
		}

		// Render name and description
		listContent += selected + base.
			Foreground(nameColor).
			Render(e.Name) + "\n"
		listContent += selected + base.
			MarginBottom(1).
			Foreground(descColor).
			Render(e.Desc) + "\n"
	}

	// Show search input if active
	var search string
	if m.isSearching {
		search = m.input.View() + "\n"
	}

	// Combine all components
	return base.
		PaddingTop(1).
		Render(title + "\n" + search + m.paginator.View() + "\n\n" + listContent)
}
