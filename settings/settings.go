package settings

import tea "github.com/charmbracelet/bubbletea"

var Model = struct {
	Name    string
	Entries []Entry
}{}

type Entry struct {
	Options  []string
	Selected int
}

func Update(msg tea.Msg) {
	switch msg.(type) {
	}
}

func View() string {
	var o string
	return o
}
