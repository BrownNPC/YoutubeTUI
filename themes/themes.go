package themes

import (
	"bytes"
	_ "embed"
	"encoding/gob"
	"fmt"
	"sync"

	"github.com/charmbracelet/lipgloss/v2"
)

type ThemeColor string

func (c ThemeColor) RGBA() (r, g, b, a uint32) {
	return lipgloss.Color(string(c)).RGBA()
}

// Theme matches the JSON schema for a terminal color theme.
type Theme struct {
	Name                string
	Black               ThemeColor
	Red                 ThemeColor
	Green               ThemeColor
	Yellow              ThemeColor
	Blue                ThemeColor
	Purple              ThemeColor
	Cyan                ThemeColor
	White               ThemeColor
	BrightBlack         ThemeColor
	BrightRed           ThemeColor
	BrightGreen         ThemeColor
	BrightYellow        ThemeColor
	BrightBlue          ThemeColor
	BrightPurple        ThemeColor
	BrightCyan          ThemeColor
	BrightWhite         ThemeColor
	Background          ThemeColor
	Foreground          ThemeColor
	CursorColor         ThemeColor
	SelectionBackground ThemeColor
}

var Themes []Theme
var ActiveID int // activeTheme index

//go:embed themes.gob
var themesGobbed []byte
var wg sync.WaitGroup

func Active() (theme Theme) {
	return Themes[ActiveID]
}
func Activate(name string) (ok bool) {
	for i, theme := range Themes {
		if string(theme.Name) == name {
			ActiveID = i
			return true
		}
	}
	return false
}
func Load() {
	wg.Add(1)
	go load()
}
func Wait() {
	wg.Wait()
}
func load() {
	defer wg.Done()
	buf := bytes.NewReader(themesGobbed)

	//Decode the GOB
	if err := gob.NewDecoder(buf).Decode(&Themes); err != nil {
		panic(fmt.Errorf("themes: gob decode failed: %w", err))
	}
}
