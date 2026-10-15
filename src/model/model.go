package model

import (
	"fmt"
	"strings"

	"github.com/aar072/mynoise-tui/classes"
	"github.com/aar072/mynoise-tui/prefs"
	"github.com/aar072/mynoise-tui/store"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	list           list.Model
	soundList      list.Model
	state          string
	detailItem     classes.Preset
	viewMode       string
	selectedCat    string
	searchInput    textinput.Model
	categories     []string
	status         string
	catalogLoading bool
	catalogError   error
	soundCache     map[string][]classes.Sound
	soundLoading   map[string]bool
	soundErrors    map[string]error
	requestedURL   string
	pendingPreset  *classes.Preset
	playbackBusy   bool
	playingURL     string
	playbackStatus string
	width, height  int
}

func NewModel() Model {
	l := list.New(nil, list.NewDefaultDelegate(), 76, 16)
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	ti := textinput.New()
	ti.Placeholder, ti.Prompt, ti.CharLimit, ti.Width = "Search...", "> ", 50, 30
	m := Model{
		list: l, state: "list", viewMode: "all", searchInput: ti,
		catalogLoading: true, soundCache: make(map[string][]classes.Sound),
		soundLoading: make(map[string]bool), soundErrors: make(map[string]error),
		width: 80, height: 24,
	}
	m.setSounds(nil)
	return m
}

func (m Model) Init() tea.Cmd { return fetchCatalogCmd }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.list.SetSize(max(1, msg.Width-4), max(1, msg.Height-7))
		m.soundList.SetSize(max(1, msg.Width-4), max(1, msg.Height-15))
		return m, nil
	case catalogMsg:
		m.catalogLoading, m.catalogError = false, msg.err
		if msg.err != nil {
			return m, nil
		}
		store.AllPresets = make(map[string]*classes.Preset)
		for _, p := range msg.presets {
			store.AllPresets[p.URL] = &classes.Preset{Data: p}
		}
		prefs.UpdatePreferences()
		m.categories = uniqueCategories(msg.presets)
		m.updateListItems()
		return m, nil
	case soundsMsg:
		delete(m.soundLoading, msg.url)
		m.soundErrors[msg.url] = msg.err
		if msg.err == nil {
			m.soundCache[msg.url] = msg.sounds
		}
		if m.detailItem.Data.URL == msg.url {
			m.setSounds(m.soundCache[msg.url])
		}
		return m, nil
	case playbackMsg:
		m.playbackBusy = false
		if msg.err != nil {
			m.playingURL = ""
			m.playbackStatus = "Playback error: " + msg.err.Error()
		} else {
			m.playingURL = msg.url
			m.playbackStatus = "Playing: " + msg.name
			if msg.generator {
				if err := prefs.RecordPlayback(msg.url); err != nil {
					m.playbackStatus += " (preferences not saved: " + err.Error() + ")"
				}
				// Refresh copied list metadata without moving the current selection.
				index := m.list.Index()
				m.updateListItems()
				m.list.Select(index)
			}
		}
		if m.detailItem.Data.URL == msg.url {
			m.status = m.playbackStatus
		}
		if pending := m.pendingPreset; pending != nil {
			m.pendingPreset = nil
			if pending.Data.URL != m.playingURL {
				m.playbackBusy = true
				m.requestedURL = pending.Data.URL
				if m.detailItem.Data.URL == pending.Data.URL {
					m.status = "Loading playback..."
				}
				return m, playPresetCmd(*pending)
			}
			if m.detailItem.Data.URL == pending.Data.URL {
				m.status = m.playbackStatus
			}
		}
		return m, nil
	}
	switch m.state {
	case "list":
		return m.handleListUpdate(msg)
	case "detail":
		return m.handleDetailUpdate(msg)
	}
	return m, nil
}

func (m *Model) setSounds(sounds []classes.Sound) {
	items := make([]list.Item, len(sounds))
	for i, sound := range sounds {
		items[i] = classes.SoundItem{Sound: sound}
	}
	m.soundList = list.New(items, list.NewDefaultDelegate(), max(1, m.width-4), max(1, m.height-15))
	m.soundList.Title = "Sounds"
	m.soundList.SetShowHelp(false)
	m.soundList.SetFilteringEnabled(false)
}

func (m *Model) startPlayback() tea.Cmd {
	if m.playbackBusy {
		m.pendingPreset = nil
		if m.detailItem.Data.URL != m.requestedURL {
			selected := m.detailItem
			m.pendingPreset = &selected
		}
		m.status = "Loading playback..."
		return nil
	}
	if m.playingURL == m.detailItem.Data.URL {
		m.status = m.playbackStatus
		return nil
	}
	m.playbackBusy = true
	m.requestedURL = m.detailItem.Data.URL
	m.status = "Loading playback..."
	return playPresetCmd(m.detailItem)
}

func (m *Model) handleItemSelection() (tea.Model, tea.Cmd) {
	if m.viewMode == "categories" {
		if selected, ok := m.list.SelectedItem().(categoryItem); ok {
			m.selectedCat, m.viewMode = selected.name, "filtered"
			m.updateListItems()
		}
		return m, nil
	}
	selected, ok := m.list.SelectedItem().(classes.Preset)
	if !ok {
		return m, nil
	}
	m.state, m.detailItem = "detail", selected
	m.setSounds(m.soundCache[selected.Data.URL])
	var fetch tea.Cmd
	if _, cached := m.soundCache[selected.Data.URL]; !cached && !m.soundLoading[selected.Data.URL] {
		m.soundLoading[selected.Data.URL] = true
		delete(m.soundErrors, selected.Data.URL)
		fetch = fetchSoundsCmd(selected.Data.URL)
	}
	return m, tea.Batch(fetch, m.startPlayback())
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
