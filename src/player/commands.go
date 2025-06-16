package player

import (
	"github.com/aar072/mynoise-tui/scraper"
	tea "github.com/charmbracelet/bubbletea"
)

// PlayPresetCmd plays the given preset (with spinner) and then emits PlaybackStartedMsg.
func PlayPresetCmd(preset scraper.Preset) tea.Cmd {
	return func() tea.Msg {
		DefaultPlayer.PlayPreset(preset)
		if DefaultPlayer.Playing {
			return PlaybackStartedMsg{PresetName: preset.Title}
		}
		// if not playing (error), emit stopped
		return PlaybackStoppedMsg{}
	}
}

// StopCmd stops playback and returns PlaybackStoppedMsg.
func StopCmd() tea.Cmd {
	return func() tea.Msg {
		DefaultPlayer.Stop()
		return PlaybackStoppedMsg{}
	}
}

// SetVolumeCmd sets the player volume and returns a VolumeChangedMsg.
func SetVolumeCmd(vol float64) tea.Cmd {
	return func() tea.Msg {
		DefaultPlayer.SetVolume(vol)
		return VolumeChangedMsg{Volume: vol}
	}
}

// SetPresenceCmd sets the noise presence and returns a PresenceChangedMsg.
func SetPresenceCmd(val float64) tea.Cmd {
	return func() tea.Msg {
		DefaultPlayer.SetPresence(val)
		return PresenceChangedMsg{Presence: val}
	}
}
