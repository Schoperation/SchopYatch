package command_tests

import (
	"errors"
	"fmt"
	"schoperation/schopyatch/command"
	"schoperation/schopyatch/msg"
	"schoperation/schopyatch/music_player"
	"testing"

	"github.com/disgoorg/disgolink/v2/lavalink"
	"github.com/stretchr/testify/require"
)

func TestSkipToCmd(t *testing.T) {
	defaultTrack := lavalink.Track{
		Encoded: "test",
		Info: lavalink.TrackInfo{
			Length: lavalink.Hour * 2,
			Author: "author",
			Title:  "title",
		},
	}

	defaultQueue := createFakeQueue(5)

	testCases := []struct {
		name                string
		inputOpts           []string
		playerConfig        playerConfig
		errorsFromPlayer    map[string]error
		expectedMessage     string
		expectedSkipMessage string
	}{
		{
			name:            "no_selection_returns_appropriate_error_message",
			expectedMessage: "No position specified. Please specify a position in the queue. E.g. `skipto 5` to go to the 5th song in the queue.",
		},
		{
			name:            "with_NaN_selection_returns_appropriate_error_message",
			inputOpts:       []string{"abc"},
			expectedMessage: "Woah hey now, that ain't a number...",
		},
		{
			name:                "no_loaded_track_returns_appropriate_error_message",
			inputOpts:           []string{"2"},
			errorsFromPlayer:    map[string]error{"SkipTo": errors.New(msg.NoLoadedTrack)},
			expectedMessage:     "Nothing to skip. Have a great evening.",
			expectedSkipMessage: fmt.Sprintf("Skipping to #%d in the queue...", 2),
		},
		{
			name:      "out_of_bounds_selection_returns_appropriate_error_message",
			inputOpts: []string{"999"},
			playerConfig: playerConfig{
				loadedTrack: &defaultTrack,
				queue:       defaultQueue,
			},
			errorsFromPlayer:    map[string]error{"SkipTo": errors.New(msg.IndexOutOfBounds)},
			expectedMessage:     fmt.Sprintf("Out of bounds. Please use a number between 1 and %d", defaultQueue.Length()),
			expectedSkipMessage: fmt.Sprintf("Skipping to #%d in the queue...", 999),
		},
		{
			name:      "end_of_queue_returns_appropriate_success_message",
			inputOpts: []string{"2"},
			playerConfig: playerConfig{
				loadedTrack: nil,
				queue:       music_player.NewMusicQueue(),
			},
			expectedMessage:     "All is now quiet on the SchopYatch front.",
			expectedSkipMessage: fmt.Sprintf("Skipping to #%d in the queue...", 2),
		},
		{
			name:      "new_track_playing_returns_appropriate_success_message",
			inputOpts: []string{"2"},
			playerConfig: playerConfig{
				loadedTrack: &defaultTrack,
				queue:       defaultQueue,
			},
			expectedMessage:     fmt.Sprintf("Now playing *%s* by **%s**.", defaultTrack.Info.Title, defaultTrack.Info.Author),
			expectedSkipMessage: fmt.Sprintf("Skipping to #%d in the queue...", 2),
		},
	}

	cmd := command.NewSkipToCmd()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc := tc
			t.Parallel()

			fakeMusicPlayer := NewDefaultFakeMusicPlayer()
			fakeMessenger := NewFakeMessenger()

			fakeMusicPlayer.ErrorsToReturn = tc.errorsFromPlayer

			fakeMusicPlayer.setPlayerConfig(tc.playerConfig)

			err := cmd.Execute(command.CommandDependencies{
				MusicPlayer: &fakeMusicPlayer,
				Messenger:   &fakeMessenger,
				Event:       NewFakeMessageCreateEvent(),
				Prefix:      ";;",
			}, tc.inputOpts...)

			require.Nil(t, err)
			require.Equal(t, tc.expectedMessage, fakeMessenger.sentMessage)
			require.Equal(t, tc.expectedSkipMessage, fakeMessenger.previousMessage)
		})
	}
}
