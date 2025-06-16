package prefs

// PresetMeta holds user preference data for a single sound preset.
type PresetMeta struct {
	IsFavorite     bool    `json:"is_favorite"`     // Whether this preset is marked as favorite
	OpenCount      int     `json:"open_count"`      // How many times this preset was opened
	FavoriteSlider float64 `json:"favorite_slider"` // User's preferred slider value for this preset
}

// UserPrefs maps preset IDs (or names) to their preference metadata.
type UserPrefs map[string]*PresetMeta
