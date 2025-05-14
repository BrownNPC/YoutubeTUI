package views

import (
	"ytt/cli"
	"ytt/components"
	"ytt/themes"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

func ChangeTheme() ChangeThemeModel {
	var rows []components.ListEntry
	for _, t := range themes.Themes {
		rows = append(rows, components.ListEntry{
			Name: string(t.Name),
			Desc: "",
		})
	}
	list := components.NewList(rows[:], "Themes")
	tabcontent := []string{
		"Themes", "Accent Color", "Selection Color",
	}
	return ChangeThemeModel{themeslist: list, tabcontent: tabcontent}
}

type ChangeThemeModel struct {
	tabcontent    []string
	selectedTab   int
	themeslist    components.List
	width, height int
}

func (m ChangeThemeModel) Update(msg tea.Msg) (ChangeThemeModel, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
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
			if active, ok := m.themeslist.Hovered(); ok {
				themes.Activate(active.Name)
				cli.Config.ThemeName = active.Name
				cli.Config.Save()
			}
		}
	}
	m.themeslist, cmd = m.themeslist.Update(msg)
	return m, cmd
}
func (m ChangeThemeModel) View() string {
	var o string
	t := themes.Active()
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
		listStyle.Render(m.themeslist.View()),
	)
	return base.Render(o)
}
