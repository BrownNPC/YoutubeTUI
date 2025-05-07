// Steal themes from https://github.com/mbadolato/iTerm2-Color-Schemes/tree/master/windowsterminal
// // github repo
// And encode it into a gob
package main

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
)

var (
	ConcurrentDownloads = 6
	SaveLocation        = "./themes/themes.gob"
)

// github repo
var (
	user   = "mbadolato"
	repo   = "iTerm2-Color-Schemes"
	branch = "master"
	folder = "windowsterminal"
)

type Theme struct {
	Name                string `json:"name"`
	Black               string `json:"black"`
	Red                 string `json:"red"`
	Green               string `json:"green"`
	Yellow              string `json:"yellow"`
	Blue                string `json:"blue"`
	Purple              string `json:"purple"`
	Cyan                string `json:"cyan"`
	White               string `json:"white"`
	BrightBlack         string `json:"brightBlack"`
	BrightRed           string `json:"brightRed"`
	BrightGreen         string `json:"brightGreen"`
	BrightYellow        string `json:"brightYellow"`
	BrightBlue          string `json:"brightBlue"`
	BrightPurple        string `json:"brightPurple"`
	BrightCyan          string `json:"brightCyan"`
	BrightWhite         string `json:"brightWhite"`
	Background          string `json:"background"`
	Foreground          string `json:"foreground"`
	CursorColor         string `json:"cursorColor"`
	SelectionBackground string `json:"selectionBackground"`
}

type GitHubContent struct {
	Name string `json:"name"`
	Type string `json:"type"` // "file" or "dir"
	Path string `json:"path"`
	URL  string `json:"download_url"` // only for files
}

var contents []GitHubContent

func init() {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s?ref=%s", user, repo, folder, branch)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching GitHub API:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("GitHub API returned status: %s\n", resp.Status)
		os.Exit(1)
	}

	if err := json.NewDecoder(resp.Body).Decode(&contents); err != nil {
		fmt.Println("Error decoding JSON:", err)
		os.Exit(1)
	}
}
func main() {

	themeChan := make(chan Theme, len(contents))
	sem := make(chan struct{}, ConcurrentDownloads) // limit to 4 concurrent downloads
	var wg sync.WaitGroup

	for _, item := range contents {
		if item.Type == "file" && path.Ext(item.Name) == ".json" {
			wg.Add(1)
			sem <- struct{}{} // acquire
			go func() {
				defer wg.Done()
				ProcessTheme(item, themeChan)
				<-sem // release
			}()
		}
	}

	// collect results in a goroutine
	go func() {
		wg.Wait()
		close(themeChan)
	}()

	var Themes []Theme
	for theme := range themeChan {
		Themes = append(Themes, theme)
	}

	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(Themes); err != nil {
		fmt.Println("Error encoding GOB:", err)
		os.Exit(1)
	}

	if err := os.WriteFile(SaveLocation, buf.Bytes(), 0644); err != nil {
		fmt.Println("Error writing file:", err)
		os.Exit(1)
	} else {
		fmt.Println("Saved themes to", SaveLocation)
	}
}

func ProcessTheme(item GitHubContent, ch chan Theme) {
	resp, err := http.Get(item.URL)
	if err != nil {
		fmt.Println("HTTP error:", err)
		return
	}
	defer resp.Body.Close()

	var theme Theme
	if err := json.NewDecoder(resp.Body).Decode(&theme); err != nil {
		fmt.Println("JSON decode error:", err)
		return
	}
	theme.Name = strings.TrimSuffix(item.Name, ".json")
	fmt.Println("Downloaded:", item.Name)
	ch <- theme
}
