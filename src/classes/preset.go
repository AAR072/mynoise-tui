package classes

type Preset struct {
	Data     ScraperPreset
	Metadata PresetMeta
}

func (p Preset) Title() string       { return p.Data.Title }
func (p Preset) Description() string { return p.Data.Category }
func (p Preset) FilterValue() string { return p.Data.Title }
