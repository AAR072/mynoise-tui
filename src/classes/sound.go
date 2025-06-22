package classes

type Sound struct {
	Sliders [10]float64
	Name    string
}

type SoundItem struct {
	Sound
}

func (s SoundItem) Title() string       { return s.Name }
func (s SoundItem) Description() string { return "" }
func (s SoundItem) FilterValue() string { return s.Name }

var DefaultSound = Sound{
	Name:    "Default Sound",
	Sliders: [10]float64{0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5},
}
