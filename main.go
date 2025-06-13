package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"

	"github.com/aar072/mynoise-tui/scraper"
)

type preset struct {
	title, url, category string
}

func (p preset) Title() string       { return p.title }
func (p preset) Description() string {
	const maxDescLen = 60
	desc := fmt.Sprintf("%s — %s", p.category, p.url)
	if len(desc) > maxDescLen {
		return desc[:maxDescLen-3] + "..."
	}
	return desc
}
func (p preset) FilterValue() string { return strings.ToLower(p.title + " " + p.category + " " + p.url) }

type model struct {
	list          list.Model
	state         string // "list" or "detail"
	detailItem    preset

	viewMode      string   // "all", "categories", "filtered"
	selectedCat   string
	searchInput   textinput.Model
	categories    []string
	allPresets    []preset
}

type categoryItem struct{ name string }

func (c categoryItem) Title() string       { return c.name }
func (c categoryItem) Description() string { return "Category" }
func (c categoryItem) FilterValue() string { return c.name }

func NewModel() model {
	presetsFromWeb, err := scraper.FetchPresets()
	if err != nil {
		fmt.Println("Error fetching presets:", err)
		os.Exit(1)
	}

	var items []list.Item
	allPresets := make([]preset, len(presetsFromWeb))
	for i, p := range presetsFromWeb {
		allPresets[i] = preset{title: p.Title, url: p.URL, category: p.Category}
		items = append(items, allPresets[i])
	}

	categories := uniqueCategories(presetsFromWeb)

	const width, height = 80, 20
	l := list.New(items, list.NewDefaultDelegate(), width, height)
	l.Title = "myNoise Presets"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	ti := textinput.New()
	ti.Placeholder = "Search..."
	ti.Prompt = "> "
	ti.CharLimit = 50
	ti.Width = 30

	m := model{
		list:        l,
		state:       "list",
		viewMode:    "all",
		categories:  categories,
		allPresets:  allPresets,
		searchInput: ti,
	}

	m.updateListItems()
	return m
}

func uniqueCategories(presets []scraper.Preset) []string {
	seen := make(map[string]struct{})
	var cats []string
	for _, p := range presets {
		if _, ok := seen[p.Category]; !ok {
			seen[p.Category] = struct{}{}
			cats = append(cats, p.Category)
		}
	}
	return cats
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.state {
	case "list":
		// Handle search input first
		if m.searchInput.Focused() {
			switch msg := msg.(type) {
			case tea.KeyMsg:
				switch msg.String() {
				case "esc":
					m.searchInput.Blur()
					m.searchInput.Reset()
					m.updateListItems()
					return m, nil
				case "enter", "down", "up":
					// Pass navigation keys to list
					m.searchInput.Blur()
					m.list, cmd = m.list.Update(msg)
					return m, cmd
				}
			}

			// Update search input for other keys
			m.searchInput, cmd = m.searchInput.Update(msg)
			m.updateListItems()
			return m, cmd
		}

		// Handle other keys when search not focused
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
			case "q", "ctrl+c":
				return m, tea.Quit
			}
		}

		// Update list for other messages
		m.list, cmd = m.list.Update(msg)
		return m, cmd

	case "detail":
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "esc", "q", "backspace":
				m.state = "list"
				return m, nil
			case "ctrl+c":
				return m, tea.Quit
			}
		}
	}

	return m, nil
}

func (m *model) handleItemSelection() (tea.Model, tea.Cmd) {
	if m.viewMode == "categories" {
		if selected, ok := m.list.SelectedItem().(categoryItem); ok {
			m.selectedCat = selected.name
			m.viewMode = "filtered"
			m.updateListItems()
		}
	} else {
		if selected, ok := m.list.SelectedItem().(preset); ok {
			m.detailItem = selected
			m.state = "detail"
		}
	}
	return m, nil
}

func (m model) View() string {
	switch m.state {
	case "list":
		header := lipgloss.NewStyle().
			Bold(true).
			MarginBottom(1).
			Render(fmt.Sprintf("Mode: %s", m.viewMode))

		if m.selectedCat != "" {
			header += "\n" + lipgloss.NewStyle().
				Foreground(lipgloss.Color("205")).
				Render("Category: " + m.selectedCat)
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
				"Title: %s\nCategory: %s\nURL: %s\n\nPress q, ESC or backspace to go back.",
				lipgloss.NewStyle().Bold(true).Render(d.title),
				lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render(d.category),
				lipgloss.NewStyle().Faint(true).Render(d.url),
			))
	}
	return ""
}

func (m *model) filterPresets() []preset {
	var filtered []preset
	q := strings.ToLower(m.searchInput.Value())

	for _, p := range m.allPresets {
		if m.selectedCat != "" && p.category != m.selectedCat {
			continue
		}

		if q != "" {
			titleMatch := strings.Contains(strings.ToLower(p.title), q)
			categoryMatch := strings.Contains(strings.ToLower(p.category), q)
			if !titleMatch && !categoryMatch {
				continue
			}
		}
		filtered = append(filtered, p)
	}
	return filtered
}

func (m *model) updateListItems() {
	var items []list.Item
	q := strings.ToLower(m.searchInput.Value())

	switch m.viewMode {
	case "all", "filtered":
		presetsToShow := m.filterPresets()
		for _, p := range presetsToShow {
			items = append(items, p)
		}
		m.list.Title = fmt.Sprintf("Presets (%d)", len(presetsToShow))

	case "categories":
		for _, c := range m.categories {
			if q == "" || strings.Contains(strings.ToLower(c), q) {
				items = append(items, categoryItem{name: c})
			}
		}
		m.list.Title = fmt.Sprintf("Categories (%d)", len(items))
	}

	m.list.SetItems(items)
	if len(items) > 0 {
		m.list.Select(0)
	}
}

func main() {
	m := NewModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
