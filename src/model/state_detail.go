package model

import (
	"time"

	"github.com/aar072/mynoise-tui/prefs"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) handleDetailUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q", "backspace":
			m.state = "list"
			return m, nil
		case "f":
			prefs.FavouritePreset(m.detailItem)
			presets := m.filterPresets()
			items := make([]list.Item, len(presets))
			for i, p := range presets {
				items[i] = p // if classes.Preset implements list.Item
				// else wrap it into a type that implements list.Item
			}

			m.list.SetItems(items)

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
