package command_tests

import (
	"fmt"
	"schoperation/schopyatch/command"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHelpCmd(t *testing.T) {
	testCases := []struct {
		name                    string
		inputOpts               []string
		expectedMessageIncludes string
	}{
		{
			name:                    "empty_opts_returns_default_message",
			expectedMessageIncludes: "Hey, SchopYatch here!\n",
		},
		{
			name:                    "invalid_opts_returns_error_message",
			inputOpts:               []string{"thiscommanddoesnotexist"},
			expectedMessageIncludes: fmt.Sprintf("Could not find %s. Try doing %shelp for a full list.", "thiscommanddoesnotexist", ";;"),
		},
		{
			name:                    "valid_opt_returns_specific_message",
			inputOpts:               []string{"help"},
			expectedMessageIncludes: "Usage: ",
		},
	}

	cmd := command.NewHelpCmd()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc := tc
			t.Parallel()

			fakeMusicPlayer := NewDefaultFakeMusicPlayer()
			fakeMessenger := NewFakeMessenger()

			err := cmd.Execute(command.CommandDependencies{
				MusicPlayer: &fakeMusicPlayer,
				Messenger:   &fakeMessenger,
				Event:       nil,
				Prefix:      ";;",
			}, tc.inputOpts...)

			require.Nil(t, err)
			require.Contains(t, fakeMessenger.sentMessage, tc.expectedMessageIncludes)
		})
	}
}
