package model

import (
	"fmt"
	"strings"
)

type Preset struct {
	TitleStr, URL, Category string
}

func (p Preset) Title() string { return p.TitleStr }
func (p Preset) Description() string {
	const maxDescLen = 60
	desc := fmt.Sprintf("%s — %s", p.Category, p.URL)
	if len(desc) > maxDescLen {
		return desc[:maxDescLen-3] + "..."
	}
	return desc
}
func (p Preset) FilterValue() string {
	return strings.ToLower(p.TitleStr + " " + p.Category + " " + p.URL)
}
