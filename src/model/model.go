package model

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/aar072/mynoise-tui/classes"
	"github.com/aar072/mynoise-tui/player"
	"github.com/aar072/mynoise-tui/prefs"
	"github.com/aar072/mynoise-tui/scraper"
	"github.com/aar072/mynoise-tui/store"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	list       list.Model
	state      string // "list", "detail", or eventually "player"
	detailItem classes.Preset

	viewMode    string // "all", "categories", "filtered"
	selectedCat string
	searchInput textinput.Model
	categories  []string
	allPresets  []classes.Preset

	status string // Current status ("Loading...", "Playing", etc.)
}

func NewModel() Model {
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

	var items []list.Item
	store.AllPresets = make([]classes.Preset, len(presetsFromWeb))
	for i, p := range presetsFromWeb {
		store.AllPresets[i] = classes.Preset{Data: p}
		items = append(items, store.AllPresets[i])
	}

	// Now we update the user preferences
	prefs.UpdatePreferences()
	prefs.SavePreferences()

	close(done)
	if err != nil {
		fmt.Println("Error fetching presets:", err)
		os.Exit(1)
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

	slices.SortFunc(store.AllPresets, func(a, b classes.Preset) int {
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
		allPresets:  store.AllPresets,
		searchInput: ti,
		status:      "",
	}

	m.updateListItems()
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle player-generated messages first
	switch msg := msg.(type) {
	case player.PlaybackStartedMsg:
		m.status = "Playing: " + msg.PresetName
		return m, nil
	case player.PlaybackStoppedMsg:
		m.status = "Stopped"
		return m, nil
	case player.VolumeChangedMsg:
		m.status = fmt.Sprintf("Volume set to %.2f", msg.Volume)
		return m, nil
	case player.PresenceChangedMsg:
		m.status = fmt.Sprintf("Presence set to %.2f", msg.Presence)
		return m, nil
	}

	// Fallback to UI state handlers
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
		if selected, ok := m.list.SelectedItem().(classes.Preset); ok {
			m.detailItem = selected
			m.state = "detail"
			m.status = "Loading..."
			return m, player.PlayPresetCmd(selected)
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
