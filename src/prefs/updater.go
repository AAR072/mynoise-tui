package prefs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aar072/mynoise-tui/classes"
	"github.com/aar072/mynoise-tui/store"
)

func UpdatePreferences() {
	if store.UserPrefs == nil {
		store.UserPrefs = make(classes.UserPrefs)
	}
	for key, preset := range store.AllPresets {
		if store.UserPrefs[key] == nil {
			store.UserPrefs[key] = &classes.PresetMeta{}
		}
		preset.Metadata = *store.UserPrefs[key]
	}
}

// SavePreferences replaces the file atomically, preserving the old file on failure.
// All preference access is owned by the Bubble Tea update loop.
func SavePreferences() error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}
	return savePreferences(path)
}

func savePreferences(path string) error {
	data, err := json.MarshalIndent(store.UserPrefs, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	file, err := os.CreateTemp(filepath.Dir(path), ".user_prefs-*")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	if _, err := file.Write(data); err != nil {
		file.Close()
		return err
	}
	if err := file.Sync(); err != nil {
		file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return os.Rename(file.Name(), path)
}

func FavouritePreset(selected *classes.Preset) error {
	preset := store.AllPresets[selected.Data.URL]
	meta := store.UserPrefs[selected.Data.URL]
	if preset == nil || meta == nil {
		return fmt.Errorf("unknown preset: %s", selected.Data.URL)
	}
	meta.IsFavorite = !meta.IsFavorite
	preset.Metadata = *meta
	selected.Metadata = *meta
	return SavePreferences()
}

func RecordPlayback(url string) error {
	preset := store.AllPresets[url]
	meta := store.UserPrefs[url]
	if preset == nil || meta == nil {
		return fmt.Errorf("unknown preset: %s", url)
	}
	meta.OpenCount++
	preset.Metadata = *meta
	return SavePreferences()
}
