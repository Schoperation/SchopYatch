package util

const (
	// Internal player
	VoiceStateNotFound = "could not find user voice state"
	NoLoadedTrack      = "player has no loaded track"
)

func IsErrorMessage(err error, errorMessage string) bool {
	return err.Error() == errorMessage
}
