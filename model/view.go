package model

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	switch m.state {
	case "list":
		header := lipgloss.NewStyle().
			Bold(true).
			MarginBottom(1).
			Render(fmt.Sprintf("Mode: %s", m.viewMode))

		if m.selectedCat != "" {
			header += "\n" + lipgloss.NewStyle().
				Foreground(lipgloss.Color("205")).
				Render("Category: "+m.selectedCat)
		}

		if m.searchInput.Focused() {
			header += "\n" + m.searchInput.View() +
				lipgloss.NewStyle().Faint(true).Render(" (press ↓/↑ to navigate, ESC to cancel)")
		} else {
			help := "Press / to search • c: categories • a: all presets • q: quit"
			header += "\n" + lipgloss.NewStyle().Faint(true).Render(help)
		}

		return lipgloss.NewStyle().
			Padding(1, 2).
			Render(header + "\n" + m.list.View())

	case "detail":
		d := m.detailItem
		return lipgloss.NewStyle().
			Margin(1, 2).
			Render(fmt.Sprintf(
				"Title: %s\nCategory: %s\nURL: %s\nStatus: %s\n\nPress q, ESC or backspace to go back.",
				lipgloss.NewStyle().Bold(true).Render(d.data.Title),
				lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render(d.data.Category),
				lipgloss.NewStyle().Faint(true).Render(d.data.URL),
				lipgloss.NewStyle().Render(m.status),
			))
	}
	return ""
}
