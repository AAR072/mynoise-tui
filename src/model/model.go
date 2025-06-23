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
	list          list.Model
	state         string // "list", "detail", or eventually "player"
	detailItem    classes.Preset
	allSounds     []classes.Sound
	selectedSound *classes.Sound
	soundList     list.Model

	viewMode    string // "all", "categories", "filtered"
	selectedCat string
	searchInput textinput.Model
	categories  []string
	allPresets  map[string]*classes.Preset

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
	store.AllPresets = make(map[string]*classes.Preset)
	for _, p := range presetsFromWeb {
		key := p.URL // or whatever unique string key exists

		presetCopy := classes.Preset{Data: p}
		store.AllPresets[key] = &presetCopy // store pointer to value

		items = append(items, presetCopy) // append value to slice
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

	m := Model{
		list:          l,
		state:         "list",
		viewMode:      "all",
		categories:    categories,
		allPresets:    store.AllPresets,
		selectedSound: &classes.DefaultSound,
		searchInput:   ti,
		status:        "",
	}

	m.updateListItems()
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle player-generated messages first
	switch msg.(type) {
	case player.PlaybackStartedMsg:
		m.status = "Playing: Default"
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
			m.state = "detail"

			// We now need to set the presets to something
			m.allSounds = scraper.FetchPresetOnclicks(selected.Data.URL)
			slices.SortFunc(m.allSounds, func(a, b classes.Sound) int {
				return strings.Compare(a.Name, b.Name)
			})
			// Now we figure out what the default sound is
			var defaultSound = scraper.GetDefaultSound(selected.Data.URL)
			m.allSounds = append([]classes.Sound{defaultSound}, m.allSounds...)
			items := make([]list.Item, len(m.allSounds))
			for i, s := range m.allSounds {
				items[i] = classes.SoundItem{s}
			}
			m.soundList = list.New(items, list.NewDefaultDelegate(), 40, 20)
			m.soundList.Title = "Presets"
			m.soundList.SetShowHelp(false)

			// If we are already playing this, we do not want to replay it
			if m.detailItem.Data.Title != selected.Title() {
				m.detailItem = selected
				return m, player.PlayPresetCmd(selected)
			}
			return m, nil
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
