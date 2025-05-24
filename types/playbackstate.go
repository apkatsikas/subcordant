package types

type PlaybackState int

const (
	AlreadyPlaying PlaybackState = iota
	PlaybackComplete
	Invalid
)
