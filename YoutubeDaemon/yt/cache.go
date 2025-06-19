package yt

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// cacheSubdir is the subdirectory under the user cache dir where playlists are stored
var cacheSubdir = filepath.Join(xdgCacheDir, "playlists")

// mutex to protect cache operations
var cacheLock sync.Mutex

// returns the playlist and true if it was successfully loaded, or false otherwise.
func loadFromCache(id string) (List, bool) {
	cacheLock.Lock()
	defer cacheLock.Unlock()

	baseCacheDir, err := os.UserCacheDir()
	if err != nil {
		// unable to determine cache directory
		return List{}, false
	}

	filePath := filepath.Join(baseCacheDir, cacheSubdir, id+".json")
	cacheFile, err := os.Open(filePath)
	if err != nil {
		// not cached or cannot open
		return List{}, false
	}
	defer cacheFile.Close()

	var p List
	if err := json.NewDecoder(cacheFile).Decode(&p); err != nil {
		// cache corrupted or invalid
		return List{}, false
	}

	return p, true
}

// saveToCache writes the given playlist to the cache.
func saveToCache(p List) {
	cacheLock.Lock()
	defer cacheLock.Unlock()

	baseCacheDir, err := os.UserCacheDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not determine cache directory:", err)
		return
	}

	cacheDir := filepath.Join(baseCacheDir, cacheSubdir)
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		fmt.Fprintln(os.Stderr, "could not create cache directory:", err)
		return
	}

	filePath := filepath.Join(cacheDir, p.ID+".json")
	cacheFile, err := os.Create(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not create cache file:", err)
		return
	}
	defer cacheFile.Close()

	if err := json.NewEncoder(cacheFile).Encode(p); err != nil {
		fmt.Fprintln(os.Stderr, "could not write playlist to cache:", err)
	}
}

// remove all cached playlist files.
// used by cli.go, the --refresh flag
func ClearCache() error {
	cacheLock.Lock()
	defer cacheLock.Unlock()

	baseCacheDir, err := os.UserCacheDir()
	if err != nil {
		return fmt.Errorf("could not determine cache directory: %w", err)
	}

	cachePath := filepath.Join(baseCacheDir, cacheSubdir)
	dirs, err := os.ReadDir(cachePath)
	if err != nil {
		return nil
	}

	for _, entry := range dirs {
		if entry.IsDir() {
			continue
		}
		if filepath.Ext(entry.Name()) == ".json" {
			if err := os.Remove(filepath.Join(cachePath, entry.Name())); err != nil {
				return fmt.Errorf("failed to delete cached file %s: %w", entry.Name(), err)
			}
		}
	}

	return nil
}
