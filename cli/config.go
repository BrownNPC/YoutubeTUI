package cli

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"ytt/themes"

	"github.com/pelletier/go-toml/v2"
)

type _config struct {
	Theme struct{
		Name string
		Accent themes.Color
	}
	Playlists []string //youtube playlist ids
}

func LoadConfig() {
	file, err := os.OpenFile(configFilePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()
	toml.NewDecoder(file).Decode(&Config)
}

// save changes
func (c _config) Save() {
	f, err := os.Create(configFilePath)
	if err != nil {
		fmt.Println(err)
	}
	err = toml.NewEncoder(f).Encode(Config)
	if err != nil {
		fmt.Println(err)
	}
	f.Close()
}
func (c *_config) AddPlaylists(ids ...string) (invalidIds []string) {
	var playlistIDRegex = regexp.MustCompile(`^PL[A-Za-z0-9_-]{32}$`)

	for _, id := range ids {
		// id must be valid
		if !playlistIDRegex.MatchString(id) {
			invalidIds = append(invalidIds, id)
			continue
		}
		// check for duplicate (already added)
		if slices.Index(c.Playlists, id) != -1 {
			continue
		}
		// finally add the id
		c.Playlists = append(c.Playlists, id)
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
