package command_tests

import (
	"errors"
	"schoperation/schopyatch/command"
	"schoperation/schopyatch/enum"
	"schoperation/schopyatch/msg"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestResumeCmd(t *testing.T) {
	testCases := []struct {
		name             string
		statusFromPlayer enum.PlayerStatus
		errorsFromPlayer map[string]error
		expectedMessage  string
	}{
		{
			name:             "no_loaded_track_returns_appropriate_error_message",
			errorsFromPlayer: map[string]error{"Unpause": errors.New(msg.NoLoadedTrack)},
			expectedMessage:  "No track's currently playing. Are we resuming a val sesh? Oh boy, it's 4 AM already, shame...",
		},
		{
			name:             "already_resumed_player_returns_appropriate_error_message",
			statusFromPlayer: enum.StatusAlreadyUnpaused,
			expectedMessage:  "Already playing. Bruh where's your ears? Better yet, where's the thing between them?",
		},
		{
			name:             "normal_circumstances_returns_success_message",
			statusFromPlayer: enum.StatusSuccess,
			expectedMessage:  "Resuming.",
		},
	}

	cmd := command.NewResumeCmd()

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
