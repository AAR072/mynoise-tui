package model

import (
	"slices"
	"strings"

	"github.com/aar072/mynoise-tui/classes"
	"github.com/aar072/mynoise-tui/player"
	"github.com/aar072/mynoise-tui/scraper"
	tea "github.com/charmbracelet/bubbletea"
)

type catalogMsg struct {
	presets []classes.ScraperPreset
	err     error
}
type soundsMsg struct {
	url    string
	sounds []classes.Sound
	err    error
}
type playbackMsg struct {
	url       string
	name      string
	generator bool
	err       error
}

func fetchCatalogCmd() tea.Msg {
	presets, err := scraper.FetchPresets()
	return catalogMsg{presets, err}
}

func fetchSoundsCmd(url string) tea.Cmd {
	return func() tea.Msg {
		sounds, err := scraper.FetchSounds(url)
		slices.SortFunc(sounds, func(a, b classes.Sound) int {
			if a.Name == b.Name {
				return 0
			}
			if a.Name == "Default" {
				return -1
			}
			if b.Name == "Default" {
				return 1
			}
			return strings.Compare(a.Name, b.Name)
		})
		return soundsMsg{url, sounds, err}
	}
}

func playPresetCmd(preset classes.Preset) tea.Cmd {
	return func() tea.Msg {
		return playbackMsg{url: preset.Data.URL, name: "Default", generator: true, err: player.PlayPreset(preset)}
	}
}

func playSoundCmd(url string, sound classes.Sound) tea.Cmd {
	return func() tea.Msg {
		return playbackMsg{url: url, name: sound.Name, err: player.PlaySound(sound)}
	}
}
