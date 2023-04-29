package command_tests

import (
	"schoperation/schopyatch/command"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPingCmd(t *testing.T) {
	testCases := []struct {
		name            string
		errorFromPlayer error
		expectedMessage string
	}{
		{
			name:            "normal_circumstances_returns_success_message",
			expectedMessage: "Pong!",
		},
	}

	cmd := command.NewPingCmd()

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
