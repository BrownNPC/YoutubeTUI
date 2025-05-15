package components

import "github.com/charmbracelet/lipgloss/v2"

func HangulFillerBorder() lipgloss.Border {
	hangulFiller := "ㅤ"
	rendered := hangulFiller
	var border = lipgloss.Border{
		Top:          rendered,
		Bottom:       rendered,
		Left:         rendered,
		Right:        rendered,
		TopLeft:      rendered,
		TopRight:     rendered,
		BottomLeft:   rendered,
		BottomRight:  rendered,
		MiddleLeft:   rendered,
		MiddleRight:  rendered,
		Middle:       rendered,
		MiddleTop:    rendered,
		MiddleBottom: rendered,
	}
	return border
}
