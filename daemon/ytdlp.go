package daemon

// abstract ytdlp cli

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

var yttDir string = getCacheDir()
var ytdlpPath = filepath.Join(getCacheDir(), ytdlpExecutableName())
var ytdlpUrl string = "https://github.com/yt-dlp/yt-dlp/releases/latest/download/" + ytdlpExecutableName()

// only function that is called from elsewhere (api.go).init()
// download ytdlp executable, dont worry about verification or overwriting
func DownloadYtdlp() error {
	// Create the file
	out, err := os.OpenFile(ytdlpPath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		if os.IsExist(err) { // already downloaded
			return nil
		}
		return errors.Join(errors.New("failed to create ytdlp file before download"), err)
	}
	defer out.Close()
	// make file writable
	err = os.Chmod(ytdlpPath, 0755)
	if err != nil {
		panic("unable to make file writable " + err.Error())
	}
	fmt.Println("downloading ytdlp, please wait")
	// Get ytdlp using http
	resp, err := http.Get(ytdlpUrl) // Get bytes
	if err != nil {
		return errors.Join(errors.New("http request failed"), err)
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	} else {
		fmt.Printf("found %s %d MB\n", ytdlpExecutableName(), resp.ContentLength/1024/1024)
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("downloaded at", ytdlpPath)
	return err
}

func runYtDLP(args ...string) (stdoutBuf, stderrBuf bytes.Buffer, err error) {
	args = append(args, "--quiet", "--no-warnings") // only errors in stderr
	cmd := exec.Command(ytdlpPath, args...)
	// Return values
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	err = cmd.Run()
	return
}

// ------------Utils------------
func ytdlpExecutableName() string {
	switch runtime.GOOS {
	case "linux":
		return "yt-dlp"
	case "windows":
		return "yt-dlp.exe"
	case "darwin":
		return "yt-dlp_macos"
	default:
		return "yt-dlp" // possibly BSD
	}
}

// get config dir path, create it if it does not exist
func getCacheDir() (dir string) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	dir = filepath.Join(home, ".cache", "ytt")

	err = os.MkdirAll(dir, 0755)
	if err != nil {
		panic(err)
	}
	return dir
}

func isYtdlpDownloaded() bool {
	_, err := os.Open(ytdlpPath)
	fmt.Println(ytdlpPath)
	return os.IsNotExist(err)
}
