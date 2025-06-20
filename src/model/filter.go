package model

import (
	"slices"
	"strings"

	"github.com/aar072/mynoise-tui/classes"
	"github.com/aar072/mynoise-tui/store"
)

func (m *Model) filterPresets() []classes.Preset {
	var filtered []classes.Preset
	q := strings.ToLower(m.searchInput.Value())

	for _, p := range store.AllPresets {
		if m.selectedCat != "" && p.Data.Category != m.selectedCat {
			continue
		}

		if q != "" {
			titleMatch := strings.Contains(strings.ToLower(p.Data.Title), q)
			categoryMatch := strings.Contains(strings.ToLower(p.Data.Category), q)
			if !titleMatch && !categoryMatch {
				continue
			}
		}
		filtered = append(filtered, *p)
	}

	slices.SortFunc(filtered, func(a, b classes.Preset) int {
		aFav := a.Metadata.IsFavorite
		bFav := b.Metadata.IsFavorite

		// Favourites go first
		if aFav && !bFav {
			return -1
		}
		if !aFav && bFav {
			return 1
		}

		// Otherwise, sort alphabetically by title
		return strings.Compare(a.Data.Title, b.Data.Title)
	})
	return filtered
}
func uniqueCategories(presets []classes.ScraperPreset) []string {
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
