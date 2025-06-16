package prefs

import (
	"encoding/json"
	"os"

	"github.com/aar072/mynoise-tui/scraper"
)

func UpdatePreferences(userPrefs UserPrefs, allPresets []scraper.Preset) UserPrefs {
	for _, presetElement := range allPresets {
		if _, exists := userPrefs[presetElement.URL]; !exists {
			userPrefs[presetElement.URL] = &PresetMeta{
				OpenCount:  0,
				IsFavorite: false,
			}
		}
	}
	return userPrefs
}

// SavePreferences saves the given UserPrefs map to the JSON config file.
func SavePreferences(userPrefs UserPrefs) error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	// Marshal userPrefs to pretty JSON
	data, err := json.MarshalIndent(userPrefs, "", "  ")
	if err != nil {
		return err
	}

	// Write JSON data to file (overwrite)
	if err := os.WriteFile(path, data, 0644); err != nil {
		return err
	}

	return nil
}
