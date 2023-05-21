package command_tests

import (
	"errors"
	"schoperation/schopyatch/command"
	"schoperation/schopyatch/enum"
	"schoperation/schopyatch/msg"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPauseCmd(t *testing.T) {
	testCases := []struct {
		name             string
		statusFromPlayer enum.PlayerStatus
		errorsFromPlayer map[string]error
		expectedMessage  string
	}{
		{
			name:             "no_loaded_track_returns_appropriate_error_message",
			errorsFromPlayer: map[string]error{"Pause": errors.New(msg.NoLoadedTrack)},
			expectedMessage:  "No track's currently playing. Would you like to pause time instead?",
		},
		{
			name:             "already_paused_player_returns_appropriate_error_message",
			statusFromPlayer: enum.StatusAlreadyPaused,
			expectedMessage:  "Already paused the music, man. Why such a party pooper?",
		},
		{
			name:             "normal_circumstances_returns_success_message",
			statusFromPlayer: enum.StatusSuccess,
			expectedMessage:  "Paused. Use `;;resume`, `;;unpause`, or `;;play` to resume.",
		},
	}

	cmd := command.NewPauseCmd()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc := tc
			t.Parallel()

			fakeMusicPlayer := NewDefaultFakeMusicPlayer()
			fakeMessenger := NewFakeMessenger()

			fakeMusicPlayer.ErrorsToReturn = tc.errorsFromPlayer
			fakeMusicPlayer.StatusToReturn = tc.statusFromPlayer

			err := cmd.Execute(command.CommandDependencies{
				MusicPlayer: &fakeMusicPlayer,
				Messenger:   &fakeMessenger,
				Event:       nil,
				Prefix:      ";;",
			})

			require.Nil(t, err)
			require.Equal(t, tc.expectedMessage, fakeMessenger.sentMessage)
		})
	}
}
