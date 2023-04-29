package command_tests

import (
	"errors"
	"schoperation/schopyatch/command"
	"schoperation/schopyatch/msg"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShuffleCmd(t *testing.T) {
	testCases := []struct {
		name            string
		errorFromPlayer error
		expectedMessage string
	}{
		{
			name:            "empty_queue_returns_appropriate_error_message",
			errorFromPlayer: errors.New(msg.QueueIsEmpty),
			expectedMessage: "Nothing to shuffle. How else am I gonna show off my riffles?",
		},
		{
			name:            "normal_circumstances_returns_success_message",
			expectedMessage: "Shuffled the queue.",
		},
	}

	cmd := command.NewShuffleCmd()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc := tc
			t.Parallel()

			fakeMusicPlayer := NewDefaultFakeMusicPlayer()
			fakeMessenger := NewFakeMessenger()

			fakeMusicPlayer.ErrorToReturn = tc.errorFromPlayer

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
