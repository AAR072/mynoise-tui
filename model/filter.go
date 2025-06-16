package model

import (
	"strings"

	"github.com/aar072/mynoise-tui/scraper"
)

func (m *Model) filterPresets() []preset {
	var filtered []preset
	q := strings.ToLower(m.searchInput.Value())

	for _, p := range m.allPresets {
		if m.selectedCat != "" && p.data.Category != m.selectedCat {
			continue
		}

		if q != "" {
			titleMatch := strings.Contains(strings.ToLower(p.data.Title), q)
			categoryMatch := strings.Contains(strings.ToLower(p.data.Category), q)
			if !titleMatch && !categoryMatch {
				continue
			}
		}
		filtered = append(filtered, p)
	}
	return filtered
}
func uniqueCategories(presets []scraper.Preset) []string {
	seen := make(map[string]struct{})
	var cats []string
	for _, p := range presets {
		if _, ok := seen[p.Category]; !ok {
			seen[p.Category] = struct{}{}
			cats = append(cats, p.Category)
		}
	}
	return cats
}
