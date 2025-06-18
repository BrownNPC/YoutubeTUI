package main

import (
	"fmt"
	"os"
	daemon "ytt/YoutubeDaemon"
	"ytt/cli"
	"ytt/themes"

	tea "github.com/charmbracelet/bubbletea/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

func main() {
	zone.NewGlobal()
	defer zone.Close()
	if cli.Run() == false {
		return
	}
	themes.Load()
	var ids []string
	for _, id := range cli.Config.Playlists {
		ids = append(ids, id)
	}
	daemon.AddPlaylists(ids...)
	themes.Wait()
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
