// Ytt cli
package cli

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"ytt/daemon"
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
	case "add", "-a":
		return AddPlaylists(args[1:])
	default:
		fmt.Println(HelpMessage)
	}
	return false
}
func RefreshCache() bool {
	dirs, err := os.ReadDir(daemon.CacheDir)
	if err != nil {
		fmt.Println(err)
		return false
	}
	for _, d := range dirs {
		if filepath.Ext(d.Name()) == ".json" {
			err = os.Remove(filepath.Join(daemon.CacheDir, d.Name()))
			if err != nil {
				fmt.Printf("err: %v\n", err)
			}

		}
	}
	return true
}
func AddPlaylists(ids []string) bool {
	if len(ids) == 0 {
		fmt.Println("Must provide YouTube playlist IDs. See: ytt help")
		return false
	}
	invalid := Config.AddPlaylists(ids...)
	for _, id := range invalid {
		fmt.Println(id, "is an invalid id.")
	}
	if len(invalid) != 0 {
		return false
	}
	Config.Save()
	return true
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
