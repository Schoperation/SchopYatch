package enum

type LoopMode int

const (
	LoopOff LoopMode = iota
	LoopTrack
	LoopQueue
)

type PlayerStatus int

const (
	StatusSuccess PlayerStatus = iota
	StatusQueued
	StatusQueuedList
	StatusPlayingAndQueuedList
	StatusAlreadyPaused
	StatusAlreadyUnpaused
	StatusFailed
)
