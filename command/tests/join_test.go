package command_tests

import (
	"errors"
	"schoperation/schopyatch/command"
	"schoperation/schopyatch/msg"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJoinCmd(t *testing.T) {
	testCases := []struct {
		name            string
		errorFromPlayer error
		expectedMessage string
	}{
		{
			name:            "no_voice_state_returns_appropriate_error_message",
			errorFromPlayer: errors.New(msg.VoiceStateNotFound),
			expectedMessage: "Dude you're not in a voice channel... get in one I can see!",
		},
		{
			name:            "normal_circumstances_returns_no_message",
			expectedMessage: "",
		},
	}

	cmd := command.NewJoinCmd()

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
				Event:       NewFakeMessageCreateEvent(),
				Prefix:      ";;",
			})

			require.Nil(t, err)
			require.Equal(t, tc.expectedMessage, fakeMessenger.sentMessage)
		})
	}
}
