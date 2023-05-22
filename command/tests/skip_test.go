package command_tests

import (
	"errors"
	"fmt"
	"schoperation/schopyatch/command"
	"schoperation/schopyatch/msg"
	"testing"

	"github.com/disgoorg/disgolink/v2/lavalink"
	"github.com/stretchr/testify/require"
)

func TestSkipCmd(t *testing.T) {
	defaultTrack := lavalink.Track{
		Encoded: "test",
		Info: lavalink.TrackInfo{
			Length: lavalink.Hour * 2,
			Author: "author",
			Title:  "title",
		},
	}

	testCases := []struct {
		name             string
		inputOpts        []string
		loadedTrack      *lavalink.Track
		errorsFromPlayer map[string]error
		expectedMessage  string
	}{
		{
			name:             "no_loaded_returns_appropriate_error_message",
			errorsFromPlayer: map[string]error{"Skip": errors.New(msg.NoLoadedTrack)},
			expectedMessage:  "Nothing to skip. Have a great evening.",
		},
		{
			name:            "end_of_queue_returns_appropriate_success_message",
			loadedTrack:     nil,
			expectedMessage: "All is now quiet on the SchopYatch front.",
		},
		{
			name:            "new_track_playing_returns_appropriate_success_message",
			loadedTrack:     &defaultTrack,
			expectedMessage: fmt.Sprintf("Now playing *%s* by **%s**.", defaultTrack.Info.Title, defaultTrack.Info.Author),
		},
	}

	cmd := command.NewSkipCmd()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc := tc
			t.Parallel()

			fakeMusicPlayer := NewDefaultFakeMusicPlayer()
			fakeMessenger := NewFakeMessenger()

			fakeMusicPlayer.ErrorsToReturn = tc.errorsFromPlayer
			fakeMusicPlayer.LoadedTrack = tc.loadedTrack

			err := cmd.Execute(command.CommandDependencies{
				MusicPlayer: &fakeMusicPlayer,
				Messenger:   &fakeMessenger,
				Event:       NewFakeMessageCreateEvent(),
				Prefix:      ";;",
			}, tc.inputOpts...)

			require.Nil(t, err)
			require.Equal(t, tc.expectedMessage, fakeMessenger.sentMessage)
			require.Equal(t, "Skipping...", fakeMessenger.previousMessage)
		})
	}
}
