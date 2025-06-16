package player

// PlaybackStartedMsg is sent when PlayPreset completes.
type PlaybackStartedMsg struct {
	PresetName string
}

// PlaybackStoppedMsg is sent when Stop completes.
type PlaybackStoppedMsg struct{}

// VolumeChangedMsg is sent when SetVolume completes.
type VolumeChangedMsg struct {
	Volume float64
}

// PresenceChangedMsg is sent when SetPresence completes.
type PresenceChangedMsg struct {
	Presence float64
}
