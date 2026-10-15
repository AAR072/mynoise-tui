package model

import (
	"fmt"
	"strconv"

	"github.com/aar072/mynoise-tui/store"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	switch m.state {
	case "list":
		if m.catalogLoading {
			return "Loading myNoise presets...\nPress ctrl+c to quit."
		}
		if m.catalogError != nil {
			return "Could not load presets: " + m.catalogError.Error() + "\nPress r to retry or ctrl+c to quit."
		}
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
			help := "Press / to search • c: categories • a: all presets"
			header += "\n" + lipgloss.NewStyle().Faint(true).Render(help)
		}

		return lipgloss.NewStyle().
			Padding(1, 2).
			Render(header + "\n" + m.list.View())

	case "detail":
		d := m.detailItem

		// Define your custom navigation instructions here, styled dim
		navInstructions := lipgloss.NewStyle().Faint(true).Render("↑/k up • ↓/j down • q back • f favourite • r retry • enter play")

		detailText := fmt.Sprintf(
			"Title: %s\nCategory: %s\nURL: %s\n%s\nListen Count: %s\nFavourited: %s\n\n%s\n\n%s",
			lipgloss.NewStyle().Bold(true).Render(d.Data.Title),
			lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render(d.Data.Category),
			lipgloss.NewStyle().Faint(true).Render(d.Data.URL),
			lipgloss.NewStyle().Render(m.status),
			lipgloss.NewStyle().Render(strconv.Itoa(store.AllPresets[d.Data.URL].Metadata.OpenCount)),
			lipgloss.NewStyle().Render(strconv.FormatBool(store.AllPresets[d.Data.URL].Metadata.IsFavorite)),
			navInstructions,
			m.soundList.View(),
		)

		if m.soundLoading[d.Data.URL] {
			detailText += "\nLoading sounds..."
		}
		if err := m.soundErrors[d.Data.URL]; err != nil {
			detailText += "\nCould not load sounds: " + err.Error() + " (r to retry)"
		}
		return lipgloss.NewStyle().
			Margin(1, 2).
			Render(detailText)
	}
	return ""
}
