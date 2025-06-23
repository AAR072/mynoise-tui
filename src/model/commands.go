package model

import (
	"github.com/aar072/mynoise-tui/browser"
	"github.com/aar072/mynoise-tui/classes"
	tea "github.com/charmbracelet/bubbletea"
)

func playPresetCmd(p classes.Preset) tea.Cmd {
	return func() tea.Msg {
		if err := browser.NavigateTo(p.Data.URL); err != nil {
			return playbackStatusMsg("Error: " + err.Error())
		}
		return checkPlaybackStatusMsg{}
	}
}
