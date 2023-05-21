package command_tests

import (
	"schoperation/schopyatch/command"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLeaveCmd(t *testing.T) {
	testCases := []struct {
		name             string
		errorsFromPlayer map[string]error
		expectedMessage  string
	}{
		{
			name:            "normal_circumstances_returns_no_message",
			expectedMessage: "",
		},
	}

	cmd := command.NewLeaveCmd()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc := tc
			t.Parallel()

			fakeMusicPlayer := NewDefaultFakeMusicPlayer()
			fakeMessenger := NewFakeMessenger()

			fakeMusicPlayer.ErrorsToReturn = tc.errorsFromPlayer

			err := cmd.Execute(command.CommandDependencies{
				MusicPlayer: &fakeMusicPlayer,
				Messenger:   &fakeMessenger,
				Event:       NewFakeMessageCreateEvent(),
				Prefix:      ";;",
			})

			require.Nil(t, err)
			require.Equal(t, tc.expectedMessage, fakeMessenger.sentMessage)
		})
	}
}
