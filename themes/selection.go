package themes

import (
	"image/color"
)

var Selection Color


func SelectionColor() color.Color {
	switch Selection {
	default:
		fallthrough
	case ThemeDefault:
		return Active().CursorColor
	case Red:
		return Active().Red
	case Green:
		return Active().Green
	case Blue:
		return Active().Blue

	case White:
		return Active().White
	case Purple:
		return Active().Purple

	case Cyan:
		return Active().Cyan
	case Yellow:
		return Active().Yellow

	// Bright Variants
	case BrightWhite:
		return Active().BrightWhite
	case BrightPurple:
		return Active().BrightPurple

	case BrightRed:
		return Active().BrightRed
	case BrightGreen:
		return Active().BrightGreen
	case BrightBlue:
		return Active().BrightBlue

	case BrightCyan:
		return Active().BrightCyan
	case BrightYellow:
		return Active().BrightYellow
	}
}
