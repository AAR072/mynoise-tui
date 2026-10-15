package prefs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aar072/mynoise-tui/classes"
	"github.com/aar072/mynoise-tui/store"
)

func TestInvalidPreferencesArePreserved(t *testing.T) {
	path := filepath.Join(t.TempDir(), "prefs.json")
	original := []byte(`{"broken":`)
	if err := os.WriteFile(path, original, 0600); err != nil {
		t.Fatal(err)
	}
	store.UserPrefs = classes.UserPrefs{"old": &classes.PresetMeta{OpenCount: 7}}
	if err := loadPreferences(path); err == nil {
		t.Fatal("expected parse error")
	}
	actual, err := os.ReadFile(path)
	if err != nil || string(actual) != string(original) {
		t.Fatal("invalid file was modified")
	}
	if store.UserPrefs["old"].OpenCount != 7 {
		t.Fatal("existing memory preferences changed")
	}
}

func TestNullPreferencesAndEntries(t *testing.T) {
	path := filepath.Join(t.TempDir(), "prefs.json")
	for _, content := range []string{"null", `{"sound":null}`} {
		if err := os.WriteFile(path, []byte(content), 0600); err != nil {
			t.Fatal(err)
		}
		if err := loadPreferences(path); err != nil {
			t.Fatal(err)
		}
		store.AllPresets = map[string]*classes.Preset{"sound": {Data: classes.ScraperPreset{URL: "sound"}}}
		UpdatePreferences()
		if store.UserPrefs["sound"] == nil {
			t.Fatal("nil metadata")
		}
	}
}

func TestAtomicPreferencesRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "new", "prefs.json")
	store.UserPrefs = classes.UserPrefs{"sound": &classes.PresetMeta{OpenCount: 3, IsFavorite: true}}
	if err := savePreferences(path); err != nil {
		t.Fatal(err)
	}
	store.UserPrefs = nil
	if err := loadPreferences(path); err != nil {
		t.Fatal(err)
	}
	if meta := store.UserPrefs["sound"]; meta.OpenCount != 3 || !meta.IsFavorite {
		t.Fatalf("bad round trip: %+v", meta)
	}
	entries, err := os.ReadDir(filepath.Dir(path))
	if err != nil || len(entries) != 1 {
		t.Fatalf("temporary files left behind: %v %v", entries, err)
	}
	// Renaming onto a directory fails; the existing destination must survive.
	if err := savePreferences(dir); err == nil {
		t.Fatal("expected rename failure")
	}
	if info, err := os.Stat(dir); err != nil || !info.IsDir() {
		t.Fatal("destination changed on failure")
	}
}
