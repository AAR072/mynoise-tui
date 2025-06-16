package model

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/aar072/mynoise-tui/prefs"
	"github.com/aar072/mynoise-tui/scraper"
)

type preset struct {
	data scraper.Preset
}

func (p preset) Title() string       { return p.data.Title }
func (p preset) Description() string { return p.data.Category }
func (p preset) FilterValue() string { return p.data.Title }

type Model struct {
	list       list.Model
	state      string // "list" or "detail"
	detailItem preset

	viewMode    string // "all", "categories", "filtered"
	selectedCat string
	searchInput textinput.Model
	categories  []string
	allPresets  []preset

	status string // Current status ("Loading...", "Playing", etc.)
}

func NewModel(userPrefs prefs.UserPrefs) Model {
	done := make(chan struct{})
	go func() {
		spinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
		i := 0
		for {
			select {
			case <-done:
				return
			default:
				fmt.Printf("\r%s Loading Sounds...", spinner[i%len(spinner)])
				i++
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()
	presetsFromWeb, err := scraper.FetchPresets()
	userPrefs = prefs.UpdatePreferences(userPrefs, presetsFromWeb)
	for url, meta := range userPrefs {
		fmt.Printf("URL: %s\n", url)
		fmt.Printf("  OpenCount: %d\n", meta.OpenCount)
		fmt.Printf("  IsFavorite: %v\n", meta.IsFavorite)
	}
	prefs.SavePreferences(userPrefs)

	// Now we update the user preferences
	close(done)
	if err != nil {
		fmt.Println("Error fetching presets:", err)
		os.Exit(1)
	}

	var items []list.Item
	allPresets := make([]preset, len(presetsFromWeb))
	for i, p := range presetsFromWeb {
		allPresets[i] = preset{data: p}
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

	slices.SortFunc(allPresets, func(a, b preset) int {
		return strings.Compare(a.Title(), b.Title())
	})
	slices.SortFunc(categories, func(a, b string) int {
		return strings.Compare(a, b)
	})

	m := Model{
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

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch m.state {
	case "list":
		return m.handleListUpdate(msg)
	case "detail":
		return m.handleDetailUpdate(msg)
	}

	return m, nil
}

func (m *Model) handleItemSelection() (tea.Model, tea.Cmd) {
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
			m.status = "Loading..."

			return m, playPresetCmd(selected)
		}
	}
	return m, nil
}

func (m *Model) updateListItems() {
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
