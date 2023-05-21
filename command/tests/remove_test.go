package command_tests

import (
	"errors"
	"fmt"
	"schoperation/schopyatch/command"
	"schoperation/schopyatch/msg"
	"schoperation/schopyatch/music_player"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRemoveCmd(t *testing.T) {
	defaultQueue := createFakeQueue(5)

	testCases := []struct {
		name             string
		inputOpts        []string
		queue            music_player.MusicQueue
		errorsFromPlayer map[string]error
		expectedMessage  string
	}{
		{
			name:            "no_selection_returns_appropriate_error_message",
			expectedMessage: fmt.Sprintf("Need to specify a number to remove. E.g. `%sremove 4` to remove the 4th track. Try `%squeue` to see the numbers.", ";;", ";;"),
		},
		{
			name:            "selection_is_NaN_returns_appropriate_error_message",
			inputOpts:       []string{"abc"},
			expectedMessage: "Dude that's some voodoo... we need a number.",
		},
		{
			name:             "empty_queue_returns_appropriate_error_message",
			inputOpts:        []string{"1"},
			queue:            music_player.NewMusicQueue(),
			errorsFromPlayer: map[string]error{"RemoveTrackFromQueue": errors.New(msg.QueueIsEmpty)},
			expectedMessage:  "Nothing to remove. The glass, this time, is fully empty.",
		},
		{
			name:             "out_of_bounds_selection_returns_appropriate_error_message",
			inputOpts:        []string{"99999"},
			queue:            defaultQueue,
			errorsFromPlayer: map[string]error{"RemoveTrackFromQueue": errors.New(msg.IndexOutOfBounds)},
			expectedMessage:  fmt.Sprintf("Out of bounds. Use a number betweem 1 and %d.", defaultQueue.Length()),
		},
		{
			name:            "valid_selection_returns_appropriate_success_message",
			inputOpts:       []string{"2"},
			queue:           defaultQueue,
			expectedMessage: fmt.Sprintf("Removed *%s* by **%s** from the queue.", defaultQueue.PeekAt(1).Info.Title, defaultQueue.PeekAt(1).Info.Author),
		},
	}

	cmd := command.NewRemoveCmd()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc := tc
			t.Parallel()

			fakeMusicPlayer := NewDefaultFakeMusicPlayer()
			fakeMessenger := NewFakeMessenger()

			fakeMusicPlayer.ErrorsToReturn = tc.errorsFromPlayer
			fakeMusicPlayer.queue = tc.queue

			err := cmd.Execute(command.CommandDependencies{
				MusicPlayer: &fakeMusicPlayer,
				Messenger:   &fakeMessenger,
				Event:       NewFakeMessageCreateEvent(),
				Prefix:      ";;",
			}, tc.inputOpts...)

			require.Nil(t, err)
			require.Equal(t, tc.expectedMessage, fakeMessenger.sentMessage)
		})
	}
}
