package components

import (
	"strings"
	"ytt/themes"

	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TableEntry struct {
	Name string
	Desc string
}

// paginated list
type List struct {
	input         textinput.Model
	isSearching   bool // search is focused
	searchChanged bool
	SearchQuery   string
	paginator     paginator.Model
	width         int
	Title         string
	AllData       []TableEntry
	FilteredData  []TableEntry
	start, end    int
	ViewHeight    int
	Cursor        int // relative to visible rows
}

// NewTable initializes a Table with entries and default viewport height.
func NewList(data []TableEntry, title string) List {
	var pag = paginator.New()
	pag.Type = paginator.Dots
	pag.ActiveDot = "◉"
	pag.SetTotalPages(len(data))
	var input = textinput.New()
	input.KeyMap.CharacterBackward.Unbind()
	input.KeyMap.CharacterForward.Unbind()
	return List{AllData: data, Title: title, paginator: pag, input: input}
}

func (m List) Update(msg tea.Msg) (List, tea.Cmd) {
	var cmd tea.Cmd
	if !m.isSearching {
		m.paginator.SetTotalPages(len(m.AllData))
		m.FilteredData = m.AllData
	} else if m.searchChanged {
		m.searchChanged = false
		var selectedData []TableEntry
		for _, e := range m.AllData {
			if m.SearchQuery != "" && !strings.Contains(strings.ToLower(e.Name+e.Desc), strings.ToLower(m.SearchQuery)) {
				continue
			}
			selectedData = append(selectedData, e)
		}
		m.paginator.SetTotalPages(len(selectedData))
		m.FilteredData = selectedData
	}
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		msg.Height = max(msg.Height, 1)
		m.ViewHeight = max((msg.Height*80)/100, 1)
		m.paginator.PerPage = max((m.ViewHeight*35)/100, 1)
	case tea.KeyMsg:
		switch msg.String() {
		case "/":
			if !m.isSearching {
				m.isSearching = true
				m.input.SetValue("")
				m.SearchQuery = ""
				cmd = m.input.Focus()
				return m, cmd
			}
		case "enter":
			if m.isSearching {
				m.isSearching = false
			}
		case "left", "h":
			m.paginator.PrevPage()
		case "right", "l":
			m.paginator.NextPage()
		case "up", "k":
			m.Cursor--
		case "down", "j":
			if m.Cursor < len(m.AllData)-1 {
				m.Cursor++
			}
		}
	}
	// Go to next page if end of current page (down)
	if m.Cursor >= m.paginator.ItemsOnPage(len(m.FilteredData)) {
		// go to next page
		if m.paginator.OnLastPage() {
			m.Cursor = m.paginator.ItemsOnPage(len(m.FilteredData)) - 1
		} else {
			m.paginator.NextPage()
			m.Cursor = 0
		}

	}
	if m.Cursor < 0 {
		if m.paginator.Page > 0 {
			m.paginator.PrevPage()
			m.Cursor = m.paginator.ItemsOnPage(len(m.FilteredData)) - 1
		} else {
			m.Cursor = 0
		}
	}
	m.Cursor %= len(m.FilteredData) + 1
	if m.isSearching {
		m.input, cmd = m.input.Update(msg)
		if m.input.Value() != m.SearchQuery {
			m.Cursor = 0
			m.paginator.Page = 0
			m.SearchQuery = m.input.Value()
			m.searchChanged = true
		}
	}
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

	m.start, m.end = m.paginator.GetSliceBounds(len(m.FilteredData))
	for i, e := range m.FilteredData[m.start:m.end] {
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
	var search string
	if m.isSearching {
		search = m.input.View() + "\n"
	}
	return base.
		PaddingTop(1).
		Render(title + "\n" + search + m.paginator.View() + "\n\n" + listContent)
}
