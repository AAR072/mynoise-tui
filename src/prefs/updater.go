package prefs

import (
	"encoding/json"
	"os"

	"github.com/aar072/mynoise-tui/classes"
	"github.com/aar072/mynoise-tui/store"
)

func UpdatePreferences() {
	for _, presetElement := range store.AllPresets {
		// If the preset does not exist
		if _, exists := store.UserPrefs[presetElement.Data.URL]; !exists {
			store.UserPrefs[presetElement.Data.URL] = &classes.PresetMeta{
				OpenCount:  0,
				IsFavorite: false,
			}
		}
		store.AllPresets[presetElement.Data.URL].Metadata = *store.UserPrefs[presetElement.Data.URL]
	}
}

// SavePreferences saves the UserPrefs map to the JSON config file.
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

func FavouritePreset(selectedPreset classes.Preset) error {
	// Flip the favourite status of the preset
	store.UserPrefs[selectedPreset.Data.URL].IsFavorite = !selectedPreset.Metadata.IsFavorite
	store.AllPresets[selectedPreset.Data.URL].Metadata.IsFavorite = !selectedPreset.Metadata.IsFavorite
	return nil
}
