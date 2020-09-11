package common

type GeneralState int

const (
	Start GeneralState = iota
	Play
	Over
)

type RelativeRectPosition struct{ XIndex, YIndex int }
