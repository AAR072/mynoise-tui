package prefs

import (
	"encoding/json"
	"os"

	"github.com/aar072/mynoise-tui/classes"
	"github.com/aar072/mynoise-tui/store"
)

func UpdatePreferences() {
	for _, presetElement := range store.AllPresets {
		if _, exists := store.UserPrefs[presetElement.Data.URL]; !exists {
			store.UserPrefs[presetElement.Data.URL] = &classes.PresetMeta{
				OpenCount:  0,
				IsFavorite: false,
			}
		} else {
			// We need to import the old preferences
			presetElement.Metadata = *store.UserPrefs[presetElement.Data.URL]
		}
	}
}

// SavePreferences saves the given UserPrefs map to the JSON config file.
func SavePreferences() error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	// Marshal store.UserPrefs to pretty JSON
	data, err := json.MarshalIndent(store.UserPrefs, "", "  ")
	if err != nil {
		return err
	}

	// Write JSON data to file (overwrite)
	if err := os.WriteFile(path, data, 0644); err != nil {
		return err
	}

	return nil
}
