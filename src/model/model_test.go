package model

import (
	"errors"
	"testing"

	"github.com/aar072/mynoise-tui/classes"
	"github.com/aar072/mynoise-tui/store"
	tea "github.com/charmbracelet/bubbletea"
)

func catalogModel(t *testing.T) Model {
	t.Helper()
	m := NewModel()
	next, _ := m.Update(catalogMsg{presets: []classes.ScraperPreset{
		{Title: "Same title", URL: "https://mynoise.net/a", Category: "Nature"},
		{Title: "Same title", URL: "https://mynoise.net/b", Category: "Nature"},
	}})
	return next.(Model)
}

func TestSelectionIsAsyncAndCached(t *testing.T) {
	m := catalogModel(t)
	next, cmd := m.handleItemSelection()
	m = *next.(*Model)
	if cmd == nil || !m.playbackBusy || m.state != "detail" {
		t.Fatal("selection did not schedule work")
	}
	url := m.detailItem.Data.URL
	next, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = *next.(*Model)
	if m.state != "list" {
		t.Fatal("navigation blocked")
	}
	next, _ = m.Update(soundsMsg{url: url, sounds: []classes.Sound{classes.DefaultSound}})
	m = next.(Model)
	if len(m.soundCache[url]) != 1 {
		t.Fatal("result not cached")
	}
	next, _ = m.Update(playbackMsg{url: url, name: "Default"})
	m = next.(Model)
	next, cmd = m.handleItemSelection()
	m = *next.(*Model)
	if cmd != nil || m.playbackBusy || m.soundLoading[url] {
		t.Fatal("cached playing preset fetched or restarted")
	}
}

func TestPlaybackFailureCanRetry(t *testing.T) {
	m := catalogModel(t)
	next, _ := m.handleItemSelection()
	m = *next.(*Model)
	next, _ = m.Update(playbackMsg{url: m.detailItem.Data.URL, err: errors.New("offline")})
	m = next.(Model)
	if m.playbackBusy || m.playingURL != "" {
		t.Fatal("failed playback marked active")
	}
	if cmd := m.startPlayback(); cmd == nil {
		t.Fatal("failed preset cannot retry")
	}
	if store.AllPresets[m.detailItem.Data.URL].Metadata.OpenCount != 0 {
		t.Fatal("failure counted as listen")
	}
}

func TestStaleResultsDoNotChangeSelectedPage(t *testing.T) {
	m := catalogModel(t)
	m.state = "detail"
	m.detailItem = *store.AllPresets["https://mynoise.net/b"]
	m.status = "Page B"
	next, _ := m.Update(soundsMsg{url: "https://mynoise.net/a", sounds: []classes.Sound{{Name: "A only"}}})
	m = next.(Model)
	if len(m.soundList.Items()) != 0 {
		t.Fatal("stale sounds displayed")
	}
	next, _ = m.Update(playbackMsg{url: "https://mynoise.net/a", name: "A only"})
	m = next.(Model)
	if m.status != "Page B" {
		t.Fatal("stale playback status displayed")
	}
	if cmd := m.startPlayback(); cmd == nil {
		t.Fatal("distinct URLs with same title treated as one preset")
	}
}

func TestCtrlCWorksDuringSearchAndLoading(t *testing.T) {
	for _, state := range []string{"list", "detail"} {
		m := NewModel()
		m.state = state
		m.searchInput.Focus()
		_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
		if cmd == nil {
			t.Fatal("missing quit command")
		}
		if _, ok := cmd().(tea.QuitMsg); !ok {
			t.Fatal("ctrl+c did not quit")
		}
	}
}

func TestRapidSelectionQueuesLatestGenerator(t *testing.T) {
	m := catalogModel(t)
	m.detailItem = *store.AllPresets["https://mynoise.net/a"]
	if cmd := m.startPlayback(); cmd == nil {
		t.Fatal("first playback not started")
	}
	m.detailItem = *store.AllPresets["https://mynoise.net/b"]
	if cmd := m.startPlayback(); cmd != nil {
		t.Fatal("concurrent playback started")
	}
	if m.pendingPreset == nil {
		t.Fatal("latest selection lost")
	}
	next, cmd := m.Update(playbackMsg{url: "https://mynoise.net/a", name: "Default"})
	m = next.(Model)
	if cmd == nil || !m.playbackBusy || m.requestedURL != "https://mynoise.net/b" {
		t.Fatal("queued generator did not start")
	}
	if m.pendingPreset != nil {
		t.Fatal("queue was not cleared")
	}
}
