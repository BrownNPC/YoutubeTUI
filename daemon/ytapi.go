package daemon

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

var CacheDir = getCacheDir()

const ytPlaylistUrl = "https://www.youtube.com/playlist?list="

type Playlist struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Channel     string  `json:"channel"`
	Entries     []Entry `json:"entries"`
}
type Entry struct {
	ID             string `json:"id"`
	URL            string `json:"url"`
	Title          string `json:"title"`
	Duration       int    `json:"duration"`
	DurationString string `json:"duration_string"`
	ChannelURL     string `json:"channel_url"`
	Uploader       string `json:"uploader"`
	ViewCount      int    `json:"view_count"`
}

func FetchPlaylist(id string) (Playlist, error) {
	if p, cached := loadFromCache(id); cached { // try loading from cache
		Playlists[id] = p
		return p, nil
	}
	stdout, stderr, err := runYtDLP(
		"--flat-playlist", "--dump-single-json", ytPlaylistUrl+id)

	if stderr.Len() != 0 { // ytdlp error
		return Playlist{}, errors.New(stderr.String())
	}
	if err != nil {
		return Playlist{}, errors.Join(errors.New("ytdlp error"), err)
	}
	var p Playlist
	err = json.Unmarshal(stdout.Bytes(), &p)
	if err != nil {
		return p, err
	}
	saveToCache(p)
	Playlists[id] = p
	return p, nil
}

// load playlist cache
func loadFromCache(id string) (p Playlist, cached bool) {
	filePath := filepath.Join(getCacheDir(), id+".json")
	cacheFile, err := os.Open(filePath)
	defer cacheFile.Close()
	if err != nil { // not cached
		return p, false
	}
	err = json.NewDecoder(cacheFile).Decode(&p)
	if err != nil {
		panic("unable to decode cached file" + err.Error())
	}
	return p, true
}

// cache playlist
func saveToCache(p Playlist) {
	filePath := filepath.Join(getCacheDir(), p.ID+".json")
	cacheFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	defer cacheFile.Close()

	if err != nil {
		if os.IsExist(err) { // already cached
			return
		}
		panic("could not write to cache " + err.Error())
	}
	err = json.NewEncoder(cacheFile).Encode(p)
	if err != nil {
		panic("could not save to cache " + err.Error())
	}
}

// audio url
type fetchedUrl struct {
	fetchedAt time.Time
	url       string
}

func newFetchedUrl(url string) fetchedUrl {
	return fetchedUrl{time.Now(), url}
}

// check if the url has expired
func (i fetchedUrl) Expired(expiry time.Duration) bool {
	return time.Since(i.fetchedAt) >= expiry
}
