package themes

import (
	"bytes"
	_ "embed"
	"encoding/gob"
	"fmt"
	"sync"

	"github.com/charmbracelet/lipgloss"
)

// Theme matches the JSON schema for a terminal color theme.
type Theme struct {
	Name                lipgloss.Color
	Black               lipgloss.Color
	Red                 lipgloss.Color
	Green               lipgloss.Color
	Yellow              lipgloss.Color
	Blue                lipgloss.Color
	Purple              lipgloss.Color
	Cyan                lipgloss.Color
	White               lipgloss.Color
	BrightBlack         lipgloss.Color
	BrightRed           lipgloss.Color
	BrightGreen         lipgloss.Color
	BrightYellow        lipgloss.Color
	BrightBlue          lipgloss.Color
	BrightPurple        lipgloss.Color
	BrightCyan          lipgloss.Color
	BrightWhite         lipgloss.Color
	Background          lipgloss.Color
	Foreground          lipgloss.Color
	CursorColor         lipgloss.Color
	SelectionBackground lipgloss.Color
}

var Themes []Theme
var activeTheme Theme

//go:embed themes.json.gob
var themesGobbed []byte
var wg sync.WaitGroup

func Active() Theme {
	return activeTheme
}
func Activate(index int) {
	activeTheme = Themes[index]
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
