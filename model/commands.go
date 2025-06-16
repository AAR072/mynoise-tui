package model

import (
	"github.com/aar072/mynoise-tui/browser"
	tea "github.com/charmbracelet/bubbletea"
)

func playPresetCmd(p preset) tea.Cmd {
	return func() tea.Msg {
		if err := browser.NavigateTo(p.data.URL); err != nil {
			return playbackStatusMsg("Error: " + err.Error())
		}
		return checkPlaybackStatusMsg{}
	}
}
func (m *Model) checkPlaybackStatus() tea.Cmd {
	return func() tea.Msg {
		if browser.IsLoading() {
			return playbackStatusMsg("Loading...")
		}
		return playbackStatusMsg("Playing")
	}
}
