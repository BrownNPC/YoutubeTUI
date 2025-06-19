package main

import (
	"fmt"
	"os"
	daemon "ytt/YoutubeDaemon"
	"ytt/YoutubeDaemon/yt"
	"ytt/cli"
	"ytt/themes"

	tea "github.com/charmbracelet/bubbletea/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

func ErrorWriter() {
	for e := range daemon.Events() {
		switch e := e.(type) {
		case daemon.EventErr:
			panic(e)
		}
	}
}

func main() {
	if cli.Run() == false {
		return
	}
	zone.NewGlobal()
	defer zone.Close()
	themes.Load()
	var ids []string
	for _, id := range cli.Config.Playlists {
		ids = append(ids, id)
	}
	<-yt.Ready
	daemon.AddPlaylists(ids...)
	themes.Wait()
	go ErrorWriter()
	themes.Activate(cli.Config.ThemeName)
	themes.Selection = cli.Config.ThemeAccent
	themes.Accent = cli.Config.ThemeAccent
	if _, err := tea.NewProgram(Model(),
		tea.WithAltScreen(),
		tea.WithMouseAllMotion(),
	).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

//   yt-dlp -f "bestaudio[ext=webm][acodec=opus]" -g
