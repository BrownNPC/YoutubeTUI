package menu

import (
	"fmt"
	"ytt/helpers"
	"ytt/themes"
	"ytt/views"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

// model
var m = struct {
	Name string
	// key to description
	Entries       []Entry
	width, height int
}{}

type Entry struct {
	Key, Description string
}

func E(key, desc string) Entry {
	return Entry{Key: key, Description: desc}
}
func init() {
	m.Entries = []Entry{
		E("p", "Go to playlist picker"),
		E("t", "Go to theme picker"),
	}
}

func Update(msg tea.Msg) (cmd tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.MouseClickMsg:
		for _, e := range m.Entries {
			if helpers.ZoneCollision(zone.Get(e.Description), msg) {
				return gotoView(e.Key)
			}
		}
	case tea.KeyMsg: // view changing time!
		return gotoView(msg.String())
	}
	return
}
func gotoView(key string) tea.Cmd {
	switch key {
	case "p":
		return views.Goto(views.ViewPlaylists)
	case "t":
		return views.Goto(views.ViewChangeTheme)
	}
	return nil
}
func View(click bool) string {
	var o string
	var t = themes.Active()
	var base = lipgloss.NewStyle().
		Foreground(t.Foreground).
		Background(t.Background)
	for i, e := range m.Entries {
		key, desc := e.Key, e.Description

		content := fmt.Sprintf("%s  %s", key, desc)
		content = zone.Mark(desc, content)
		var newline string
		if i != len(m.Entries)-1 {
			newline = "\n"
		}
		o += base.
			PaddingLeft(1).
			PaddingRight(1).
			Render(content) + newline
	}
	border := lipgloss.NormalBorder()
	o = base.
		Border(border).
		BorderForeground(t.Foreground).
		BorderBackground(t.Background).
		Render(o)

	space := base.
		Render("Space")
	if click {
		space = base.Render("Right Click")
	}
	o, _ = helpers.Overlay(o, space, 0, 1, false)
	return o
}
