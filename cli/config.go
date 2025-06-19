package cli

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"
	"ytt/themes"

	"github.com/pelletier/go-toml/v2"
)

type _config struct {
	ThemeName           string
	ThemeAccent         themes.Color
	ThemeSelectionColor themes.Color
	Playlists           []string //youtube playlist ids
}

func LoadConfig() {
	file, err := os.OpenFile(configFilePath, os.O_RDWR|os.O_CREATE, 0644)
	defer file.Close()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	toml.NewDecoder(file).Decode(&Config)
}

// save changes
func (c _config) Save() {
	f, err := os.Create(configFilePath)
	defer f.Close()
	if err != nil {
		fmt.Println(err)
	}
	err = toml.NewEncoder(f).Encode(Config)
	if err != nil {
		panic(err)
	}
}

// https://stackoverflow.com/a/75373610
var playlistIDRegex = regexp.MustCompile(`[?&]list=([^#?&]*)`)

func (c *_config) AddPlaylists(inputs ...string) (invalidIds []string) {
	for _, input := range inputs {
		input = strings.TrimSpace(input)

		var id string
		if m := playlistIDRegex.FindStringSubmatch(input); m != nil {
			id = m[1]
		} else {
			invalidIds = append(invalidIds, input)
			continue
		}

		if slices.Index(c.Playlists, id) == -1 {
			c.Playlists = append(c.Playlists, id)
		}
	}
	return invalidIds
}
func (c *_config) RemovePlaylist(playlistId string) {
	i := slices.Index(c.Playlists, playlistId)
	if i == -1 {
		return
	}
	c.Playlists = slices.Delete(c.Playlists, i, i)
}
