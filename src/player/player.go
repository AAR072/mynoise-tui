package player

import (
	"fmt"

	"github.com/aar072/mynoise-tui/browser"
	"github.com/aar072/mynoise-tui/classes"
)

// Player holds the playback state and slider values.
type Player struct {
	CurrentPreset classes.Preset
	Volume        float64
	Presence      float64
	Playing       bool
}

// New creates a fresh Player with default slider settings.
func New() *Player {
	return &Player{
		Volume:   1.0,
		Presence: 1.0,
		Playing:  false,
	}
}

// PlayPreset begins playback of the given preset.
// It blocks until `browser.Play` returns.
func (p *Player) PlayPreset(preset classes.Preset) {
	// We open the browser
	err := browser.NavigateTo(preset.Data.URL)
	// We increment the opencount
	preset.Metadata.OpenCount++

	if err != nil {
		fmt.Println("Error playing preset:", err)
		p.Playing = false
		return
	}

	p.CurrentPreset = preset
	p.Playing = true
}

// Stop ends playback.
func (p *Player) Stop() {
	p.Playing = false
	// TODO: send a stop signal to your audio backend here
}

// SetVolume adjusts the playback volume (0.0–1.0).
func (p *Player) SetVolume(volume float64) {
	p.Volume = volume
	// TODO: propagate this change to your audio backend
}

// SetPresence adjusts the noise “presence” slider (0.0–1.0).
func (p *Player) SetPresence(presence float64) {
	p.Presence = presence
	// TODO: propagate this change to your audio backend
}

// DefaultPlayer is the package‑level singleton you can use directly.
var DefaultPlayer = New()
