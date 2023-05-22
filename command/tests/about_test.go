package command_tests

import (
	"fmt"
	"schoperation/schopyatch/bot"
	"schoperation/schopyatch/command"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAboutCmd(t *testing.T) {
	testCases := []struct {
		name             string
		errorsFromPlayer map[string]error
		expectedMessage  string
	}{
		{
			name: "normal_circumstances_returns_success_message",
			expectedMessage: "```" +
				fmt.Sprintf("SchopYatch v%s\n\n", bot.SchopYatchVersion) +
				"Coded by Schoperation: 		   https://github.com/Schoperation/SchopYatch\n" +
				"Lavalink Client by Freya Arbjerg: https://github.com/freyacodes/Lavalink-Client\n" +
				"Libraries written by the DisGoOrg:\n" +
				"\tDisGo:     https://github.com/DisgoOrg/disgo\n" +
				"\tDisGoLink: https://github.com/disgoorg/disgolink\n" +
				"```",
		},
	}

	cmd := command.NewAboutCmd()

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
				Event:       nil,
				Prefix:      ";;",
				Version:     bot.SchopYatchVersion,
			})

			require.Nil(t, err)
			require.Equal(t, tc.expectedMessage, fakeMessenger.sentMessage)
		})
	}
}
