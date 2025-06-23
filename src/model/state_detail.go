package model

import (
	"github.com/aar072/mynoise-tui/classes"
	"github.com/aar072/mynoise-tui/player"
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
			prefs.FavouritePreset(&m.detailItem)
			presets := m.filterPresets()
			items := make([]list.Item, len(presets))
			for i, p := range presets {
				items[i] = p
			}
			m.list.SetItems(items)
			return m, nil

		case "enter":
			if selected, ok := m.soundList.SelectedItem().(classes.SoundItem); ok {
				m.status = "Playing: " + selected.Name

				player.DefaultPlayer.PlaySound(selected.Name, selected.Sliders)
			}
			return m, nil

		case "ctrl+c":
			return m, tea.Quit
		}

	}

	// Handle navigation in soundList
	var cmd tea.Cmd
	m.soundList, cmd = m.soundList.Update(msg)
	return m, cmd
}
