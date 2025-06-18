package model

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) handleListUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.searchInput.Focused() {
		switch msg := msg.(type) {
		case playbackStatusMsg:
			m.status = string(msg)
			return m, nil
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				m.searchInput.Blur()
				m.searchInput.Reset()
				m.updateListItems()
				return m, nil
			case "enter", "down", "up":
				m.searchInput.Blur()
				m.list, cmd = m.list.Update(msg)
				return m, cmd
			}
		}

		m.searchInput, cmd = m.searchInput.Update(msg)
		m.updateListItems()
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m.handleItemSelection()
		case "/":
			m.searchInput.Focus()
			return m, nil
		case "esc":
			if m.selectedCat != "" {
				m.selectedCat = ""
				m.viewMode = "all"
				m.updateListItems()
			}
			return m, nil
		case "c":
			m.viewMode = "categories"
			m.updateListItems()
			return m, nil
		case "a":
			m.viewMode = "all"
			m.selectedCat = ""
			m.updateListItems()
			return m, nil
		}
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}
