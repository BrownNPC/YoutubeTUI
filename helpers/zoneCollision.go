package helpers

import (
	tea "github.com/charmbracelet/bubbletea/v2"
	zone "github.com/lrstanley/bubblezone/v2"
)

func ZoneCollision(z *zone.ZoneInfo, msg tea.MouseMsg) bool {
	if z.IsZero() {
		return false
	}
	mouse := msg.Mouse()
	return mouse.X >= z.StartX && mouse.X <= z.EndX &&
		mouse.Y >= z.StartY && mouse.Y <= z.EndY
}
