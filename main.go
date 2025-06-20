package main

import (
	"fmt"
	"os"
	"time"
	daemon "ytt/YoutubeDaemon"
	"ytt/cli"
	"ytt/themes"

	tea "github.com/charmbracelet/bubbletea/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

func ErrorWriter() {
	logfile, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Errorf("Coult not create log file %w", err))
	}
	for e := range daemon.Events() {
		switch e := e.(type) {
		case daemon.EventErr:
			fmt.Fprintln(logfile, time.Now(), e)
		case daemon.EventInfo:
			fmt.Fprintln(logfile, time.Now(), e)
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
	daemon.InitDaemon()
	daemon.RegisterPlaylists(ids...)
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
