package prefs

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/aar072/mynoise-tui/classes"
	"github.com/aar072/mynoise-tui/store"
)

// getConfigPath returns the full path to the JSON config file storing user prefs.
// It uses the OS-specific user config directory and appends "mynoise-tui/user_prefs.json".
func getConfigPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "mynoise-tui", "user_prefs.json"), nil
}

// InitConfig ensures the user prefs JSON file exists, loading it if present,
// or creating an empty one otherwise. Returns the loaded preferences map.
func InitConfig() (classes.UserPrefs, error) {
	if store.UserPrefs == nil {
		store.UserPrefs = make(map[string]*classes.PresetMeta)
	}
	path, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	// Check if the prefs file exists.
	_, err = os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		// Create the directory structure if it doesn't exist.
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return nil, err
		}
		// Initialize an empty prefs map.
		emptyPrefs := classes.UserPrefs{}
		// Marshal empty prefs to pretty JSON.
		data, err := json.MarshalIndent(emptyPrefs, "", "  ")
		if err != nil {
			return nil, err
		}
		// Write the empty JSON file to disk.
		if err := os.WriteFile(path, data, 0644); err != nil {
			return nil, err
		}
	} else if err != nil {
		// Some other error accessing the file.
		return nil, err
	}

	// Read the JSON prefs file.
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var prefs classes.UserPrefs
	// Unmarshal JSON data into the prefs map.
	if err := json.Unmarshal(data, &prefs); err != nil {
		return nil, err
	}

	return prefs, nil
}
