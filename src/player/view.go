package player

import (
	"fmt"
)

// View returns a simple text summary of the player's current state.
// You can call this from your model's View() when in the "player" state.
func View() string {
	if DefaultPlayer.Playing {
		return fmt.Sprintf(
			"▶ %s\nVolume: %.2f | Presence: %.2f",
			DefaultPlayer.CurrentPreset.Title,
			DefaultPlayer.Volume,
			DefaultPlayer.Presence,
		)
	}
	return "■ Stopped"
}
