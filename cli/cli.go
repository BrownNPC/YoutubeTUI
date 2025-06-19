// Ytt cli
package cli

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"ytt/YoutubeDaemon/yt"
)

var configDir = func() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".config", "ytt")
}()

//go:embed help.txt
var HelpMessage string
var (
	configFilePath = filepath.Join(configDir, "config.toml")
	Config         _config
)

func init() {
	// make sure config directory exists
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		if !os.IsExist(err) {
			fmt.Println("error while creating config directory: ", err)
		}
	}
	LoadConfig()
}

// returning false means we should quit after Run
func Run() (run bool) {
	if len(os.Args) >= 2 {
		return HandleArgs(os.Args[1:]...)
	}
	return true
}

func HandleArgs(args ...string) (run bool) {
	switch args[0] {
	case "help", "-h":
		fmt.Println(HelpMessage)
	case "refresh", "-r":
		return RefreshCache()
	case "config", "-c":
		OpenConfigDir()
		return false
	case "add", "-a":
		return AddPlaylists(args[1:])
	default:
		fmt.Println(HelpMessage)
	}
	return false
}
func RefreshCache() (run bool) {
	err := yt.ClearCache()
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
func AddPlaylists(urls []string) bool {
	if len(urls) == 0 {
		fmt.Println("Must provide YouTube playlist URL.")
		fmt.Println("Example: ", `ytt add "https://www.youtube.com/watch?v=0QvdDX2Q7rI&list=PLN1mxegxWPd0GfRvWy_WzwpNKnqSWTV5U"`)
		return false
	}
	invalid := Config.AddPlaylists(urls...)
	for _, id := range invalid {
		fmt.Println(id, "is an invalid url, please make sure to surround the url with double quotes.")
	}
	if len(invalid) != 0 {
	}
	Config.Save()
	fmt.Println("Added!")
	return false
}
func OpenConfigDir() {
	var cmd *exec.Cmd
	path := configDir
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("explorer", path)
	case "darwin":
		cmd = exec.Command("open", path)
	case "linux":
		cmd = exec.Command("xdg-open", path)
	default:
		cmd = exec.Command("xdg-open", path)
	}

	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
	}
}
