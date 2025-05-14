package views

import tea "github.com/charmbracelet/bubbletea/v2"

type ViewMsg int

const (
	ViewPlaylists ViewMsg = iota
	ViewChangeTheme
)

func Goto(v ViewMsg) tea.Cmd {
	return func() tea.Msg {
		return v
	}
}
