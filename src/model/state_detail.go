package model

import (
	"github.com/aar072/mynoise-tui/classes"
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
			if err := prefs.FavouritePreset(&m.detailItem); err != nil {
				m.status = "Preferences not saved: " + err.Error()
			}
			presets := m.filterPresets()
			items := make([]list.Item, len(presets))
			for i, p := range presets {
				items[i] = p
			}
			m.list.SetItems(items)
			return m, nil

		case "r":
			url := m.detailItem.Data.URL
			if !m.soundLoading[url] {
				m.soundLoading[url] = true
				delete(m.soundErrors, url)
				return m, tea.Batch(fetchSoundsCmd(url), m.startPlayback())
			}
			return m, m.startPlayback()
		case "enter":
			if m.playbackBusy {
				return m, nil
			}
			if m.playingURL != m.detailItem.Data.URL {
				return m, m.startPlayback()
			}
			if selected, ok := m.soundList.SelectedItem().(classes.SoundItem); ok {
				m.playbackBusy = true
				m.requestedURL = m.detailItem.Data.URL
				m.status = "Loading: " + selected.Name
				return m, playSoundCmd(m.detailItem.Data.URL, selected.Sound)
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
