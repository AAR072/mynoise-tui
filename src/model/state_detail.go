package model

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) handleDetailUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "backspace":
			m.state = "list"
			return m, nil
		case "ctrl+c":
			return m, tea.Quit
		}

	case playbackStatusMsg:
		m.status = string(msg)
		// Continue checking status until we get "Playing"
		if m.status != "Playing" {
			return m, tea.Tick(500*time.Millisecond, func(time.Time) tea.Msg {
				return checkPlaybackStatusMsg{}
			})
		}
		return m, nil

	case checkPlaybackStatusMsg:
		return m, m.checkPlaybackStatus()
	}

	return m, nil
}
