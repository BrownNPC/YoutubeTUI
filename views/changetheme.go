package views

import (
	"ytt/cli"
	"ytt/components"
	"ytt/helpers"
	"ytt/themes"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

func ChangeTheme() ChangeThemeModel {
	// Themes
	var rows []components.ListEntry
	for _, t := range themes.Themes {
		rows = append(rows, components.ListEntry{
			Name: string(t.Name),
			Desc: "",
		})
	}
	themelist := components.NewList(rows, "Themes")
	//Accents
	rows = []components.ListEntry{
		components.ListEntry{Name: "ThemeDefault"},
		components.ListEntry{Name: "Red"},
		components.ListEntry{Name: "Green"},
		components.ListEntry{Name: "Blue"},
		components.ListEntry{Name: "White"},
		components.ListEntry{Name: "Purple"},
		components.ListEntry{Name: "Yellow"},
		components.ListEntry{Name: "Pink"},
		components.ListEntry{Name: "Cyan"},
		components.ListEntry{Name: "BrightWhite"},
		components.ListEntry{Name: "BrightPurple"},
		components.ListEntry{Name: "BrightRed"},
		components.ListEntry{Name: "BrightGreen"},
		components.ListEntry{Name: "BrightBlue"},
		components.ListEntry{Name: "BrightYellow"},
		components.ListEntry{Name: "BrightCyan"},
	}

	accentList := components.NewList(rows, "Accent Colors")
	selectionColorList := components.NewList(rows, "Selection Colors")

	// blink the currently activew themes
	themelist.SelectedName = string(themes.Active().Name)
	accentList.SelectedName = string(themes.Accent)
	selectionColorList.SelectedName = string(themes.Selection)

	tabcontent := []string{
		"Themes", "Accent Color", "Selection Color",
	}
	m := ChangeThemeModel{
		themeslist:    themelist,
		tabcontent:    tabcontent,
		accentsList:   accentList,
		selectionList: selectionColorList,
	}
	return m
}

type ChangeThemeModel struct {
	themeslist    components.List
	accentsList   components.List
	selectionList components.List
	tabcontent    []string
	selectedTab   int
	width, height int
}

func (m ChangeThemeModel) Update(msg tea.Msg) (ChangeThemeModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.themeslist, cmd = m.themeslist.Update(msg)
		m.accentsList, cmd = m.accentsList.Update(msg)
		m.selectionList, cmd = m.selectionList.Update(msg)
	case tea.MouseMsg:
		for i, c := range m.tabcontent {
			z := zone.Get(c)
			if helpers.ZoneCollision(z, msg) && msg.Mouse().Button == tea.MouseLeft {
				m.selectedTab = i
			}
		}
		// List content
		switch m.tabcontent[m.selectedTab] {
		case "Themes":
			if active, ok := m.themeslist.MouseHovered(msg); ok {
				if msg.Mouse().Button == tea.MouseLeft {
					m.updateTheme(active)
				}
			}
		case "Accent Color":
			if active, ok := m.accentsList.MouseHovered(msg); ok {
				if msg.Mouse().Button == tea.MouseLeft {
					m.updateAccent(active)
				}
			}
		case "Selection Color":
			if active, ok := m.selectionList.MouseHovered(msg); ok {
				if msg.Mouse().Button == tea.MouseLeft {
					m.updateSelectedColor(active)
				}
			}
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			m.selectedTab = (m.selectedTab + 1) % len(m.tabcontent)
		case "shift+tab":
			m.selectedTab--
			if m.selectedTab < 0 {
				m.selectedTab = len(m.tabcontent) - 1
			}
		case "enter":
			switch m.tabcontent[m.selectedTab] {
			case "Themes":
				if active, ok := m.themeslist.Hovered(); ok {
					m.updateTheme(active)
				}
			case "Accent Color":
				if active, ok := m.accentsList.Hovered(); ok {
					m.updateAccent(active)
				}
			case "Selection Color":
				if active, ok := m.selectionList.Hovered(); ok {
					m.updateSelectedColor(active)
				}
			}
		}
	}
	switch m.tabcontent[m.selectedTab] {
	case "Themes":
		m.themeslist, cmd = m.themeslist.Update(msg)
	case "Accent Color":
		m.accentsList, cmd = m.accentsList.Update(msg)
	case "Selection Color":
		m.selectionList, cmd = m.selectionList.Update(msg)
	}
	return m, cmd
}
func (m ChangeThemeModel) View() string {
	var o string
	t := themes.Active()
	var visibleList components.List
	switch m.tabcontent[m.selectedTab] {
	case "Themes":
		visibleList = m.themeslist
	case "Accent Color":
		visibleList = m.accentsList
	case "Selection Color":
		visibleList = m.selectionList
	}

	var base = lipgloss.NewStyle().
		Background(t.Background)
	var tabStyle = base.
		Border(lipgloss.NormalBorder()).
		BorderBackground(t.Background)
	var tabs []string
	for i, c := range m.tabcontent {
		var content string
		if i == m.selectedTab {
			content = tabStyle.
				BorderForeground(themes.AccentColor()).
				Foreground(themes.SelectionColor()).
				Render(c)
		} else {
			content = tabStyle.
				Render(c)
		}
		content = zone.Mark(c, content)
		tabs = append(tabs, content)
	}
	tabContent := lipgloss.JoinHorizontal(0, tabs...)
	tabContent = base.
		Width(m.width).
		Render(tabContent)
	listStyle := base.
		Width(m.width).
		Height(m.height).
		PaddingLeft(2).
		MarginTop(0)
	o += lipgloss.JoinVertical(0,
		tabContent,
		listStyle.Render(visibleList.View()),
	)
	return base.Render(o)
}
func (m *ChangeThemeModel) updateTheme(active components.ListEntry) {
	m.themeslist.SelectedName = active.Name
	themes.Activate(active.Name)
	cli.Config.ThemeName = active.Name
	cli.Config.Save()
}
func (m *ChangeThemeModel) updateAccent(active components.ListEntry) {
	m.accentsList.SelectedName = active.Name
	themes.Accent = themes.Color(active.Name)
	cli.Config.ThemeAccent = themes.Color(active.Name)
	cli.Config.Save()
}
func (m *ChangeThemeModel) updateSelectedColor(active components.ListEntry) {
	m.selectionList.SelectedName = active.Name
	themes.Selection = themes.Color(active.Name)
	cli.Config.ThemeSelectionColor = themes.Color(active.Name)
	cli.Config.Save()
}
