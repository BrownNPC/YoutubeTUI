package main

import (
	"fmt"
	"github.com/charmbracelet/bubbletea/v2"
	"os"
	"sync"
	"ytt/cli"
	"ytt/daemon"
	"ytt/themes"
)

func main() {

	if cli.Run() == false {
		return
	}
	themes.Load()
	var wg sync.WaitGroup
	for _, id := range cli.Config.Playlists {
		wg.Add(1)
		go fillCache(&wg, id)
	}
	wg.Wait()
	themes.Wait()
	if _, err := tea.NewProgram(Model(),
		tea.WithAltScreen(),
		// tea.WithMouseCellMotion(),
	).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func fillCache(wg *sync.WaitGroup, id string) {
	defer wg.Done()
	_, err := daemon.FetchPlaylist(id)
	if err != nil {
		fmt.Println(err)
	}
}

//
//   yt-dlp -f "bestaudio[ext=webm][acodec=opus]" -g
