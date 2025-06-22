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

func FavouritePreset(selectedPreset *classes.Preset) error {
	// Flip the favorite status of the preset in the global store
	newFav := !selectedPreset.Metadata.IsFavorite
	store.UserPrefs[selectedPreset.Data.URL].IsFavorite = newFav
	store.AllPresets[selectedPreset.Data.URL].Metadata.IsFavorite = newFav

	// Update the selectedPreset itself so caller sees the change
	selectedPreset.Metadata.IsFavorite = newFav
	return nil
}
