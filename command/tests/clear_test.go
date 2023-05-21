package command_tests

import (
	"fmt"
	"schoperation/schopyatch/command"
	"testing"

	"github.com/disgoorg/disgolink/v2/lavalink"
	"github.com/stretchr/testify/require"
)

func TestClearCmd(t *testing.T) {
	testCases := []struct {
		name             string
		isFailureTest    bool
		inputOpts        []string
		isQueueEmpty     bool
		errorsFromPlayer map[string]error
		expectedMessage  string
	}{
		{
			name:            "empty_queue_returns_appropriate_message",
			isQueueEmpty:    true,
			expectedMessage: "Already clear. Were you hoping for a funny error? Same",
		},
		{
			name:            "invalid_option_returns_appropriate_error_message",
			isFailureTest:   true,
			inputOpts:       []string{"e"},
			expectedMessage: "Dude, that was not a number...",
		},
		{
			name:            "valid_option_returns_success_message",
			inputOpts:       []string{"2"},
			expectedMessage: fmt.Sprintf("Cleared the first %d tracks from the queue.", 2),
		},
		{
			name:            "no_option_returns_success_message",
			expectedMessage: "Cleared the queue.",
		},
	}

	cmd := command.NewClearCmd()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc := tc
			t.Parallel()

			fakeMusicPlayer := NewDefaultFakeMusicPlayer()
			fakeMessenger := NewFakeMessenger()

			fakeMusicPlayer.ErrorsToReturn = tc.errorsFromPlayer

			if !tc.isQueueEmpty {
				fakeMusicPlayer.AddTrackToQueue(lavalink.Track{})
			}

			err := cmd.Execute(command.CommandDependencies{
				MusicPlayer: &fakeMusicPlayer,
				Messenger:   &fakeMessenger,
				Event:       nil,
				Prefix:      ";;",
			}, tc.inputOpts...)

			require.Nil(t, err)
			require.Equal(t, tc.expectedMessage, fakeMessenger.sentMessage)

			if !tc.isFailureTest {
				require.True(t, fakeMusicPlayer.IsQueueEmpty())
			}
		})
	}
}
