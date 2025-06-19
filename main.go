package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	daemon "ytt/YoutubeDaemon"
	"ytt/YoutubeDaemon/yt"
	"ytt/cli"
	"ytt/themes"

	tea "github.com/charmbracelet/bubbletea/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

func GetLatestYtDlpVersion() (string, error) {
	resp, err := http.Get("https://api.github.com/repos/yt-dlp/yt-dlp/releases/latest")
	if err != nil {
		return "", fmt.Errorf("failed to fetch release info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("non-200 response from GitHub: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var data struct {
		TagName string `json:"tag_name"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	return data.TagName, nil
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
