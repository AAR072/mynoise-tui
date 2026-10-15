package prefs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aar072/mynoise-tui/classes"
	"github.com/aar072/mynoise-tui/store"
)

func getConfigPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "mynoise-tui", "user_prefs.json"), nil
}

// InitConfig never overwrites unreadable or invalid existing preferences.
func InitConfig() error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}
	return loadPreferences(path)
}

func loadPreferences(path string) error {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		store.UserPrefs = make(classes.UserPrefs)
		return nil
	}
	if err != nil {
		return fmt.Errorf("read preferences: %w", err)
	}
	var loaded classes.UserPrefs
	if err := json.Unmarshal(data, &loaded); err != nil {
		return fmt.Errorf("parse preferences %s: %w", path, err)
	}
	if loaded == nil {
		loaded = make(classes.UserPrefs)
	}
	for key, meta := range loaded {
		if meta == nil {
			loaded[key] = &classes.PresetMeta{}
		}
	}
	store.UserPrefs = loaded
	return nil
}
