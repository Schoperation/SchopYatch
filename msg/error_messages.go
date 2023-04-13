package msg

const (
	// Internal player
	VoiceStateNotFound = "could not find user voice state"
	NoLoadedTrack      = "player has no loaded track"
	AlreadyLoadedTrack = "player already has a loaded track"
	QueueIsEmpty       = "queue is empty"
	NoResultsFound     = "no results found"
	IndexOutOfBounds   = "index out of bounds"
	NegativeDuration   = "duration cannot be negative"
)

func IsErrorMessage(err error, errorMessage string) bool {
	return err.Error() == errorMessage
}
