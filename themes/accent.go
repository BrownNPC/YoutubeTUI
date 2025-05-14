package themes

import (
	"image/color"

)

type Color string

const (
	Red   Color = "Red"
	Green       = "Green"
	Blue        = "Blue"

	White  = "White"
	Purple = "Purple"

	Yellow = "Yellow"
	Pink   = "Pink"
	Cyan   = "Cyan"

	BrightWhite  = "BrightWhite"
	BrightPurple = "BrightPurple"

	BrightRed   = "BrightRed"
	BrightGreen = "BrightGreen"
	BrightBlue  = "BrightBlue"

	BrightYellow = "BrightYellow"
	BrightCyan   = "BrightCyan"
)

var ActiveAccent Color

func SetAccent(a Color) {
	ActiveAccent = a
}
func AccentColor() color.Color {
	switch ActiveAccent {
	default:
		fallthrough
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
