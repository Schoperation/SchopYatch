package command_tests

import (
	"errors"
	"fmt"
	"schoperation/schopyatch/command"
	"schoperation/schopyatch/enum"
	"schoperation/schopyatch/msg"
	"testing"

	"github.com/disgoorg/disgolink/v2/lavalink"
	"github.com/stretchr/testify/require"
)

func TestNowPlayingCmd(t *testing.T) {
	defaultTrack := lavalink.Track{
		Encoded: "test",
		Info: lavalink.TrackInfo{
			Length: lavalink.Hour,
			Author: "author",
			Title:  "title",
		},
	}

	testCases := []struct {
		name             string
		errorsFromPlayer map[string]error
		loopMode         enum.LoopMode
		track            lavalink.Track
		position         lavalink.Duration
		expectedMessage  string
	}{
		{
			name:             "no_loaded_track_returns_appropriate_error_message",
			errorsFromPlayer: map[string]error{"GetLoadedTrack": errors.New(msg.NoLoadedTrack)},
			expectedMessage:  "Nothing's playing. Bruh moment...",
		},
		{
			name:     "normal_circumstances_returns_appropriate_message",
			loopMode: enum.LoopOff,
			track:    defaultTrack,
			position: lavalink.Minute,
			expectedMessage: fmt.Sprintf("Now Playing:\n\t*%s* by **%s**\n\t%s `[%s / %s]`\n\t%s",
				defaultTrack.Info.Title,
				defaultTrack.Info.Author,
				"59 mins, 0 secs remaining.",
				lavalink.Minute,
				defaultTrack.Info.Length.String(),
				""),
		},
		{
			name:     "with_complex_position_returns_appropriate_message",
			loopMode: enum.LoopOff,
			track:    defaultTrack,
			position: lavalink.Duration(125000), // 2 minutes, 5 seconds
			expectedMessage: fmt.Sprintf("Now Playing:\n\t*%s* by **%s**\n\t%s `[%s / %s]`\n\t%s",
				defaultTrack.Info.Title,
				defaultTrack.Info.Author,
				"57 mins, 55 secs remaining.",
				lavalink.Duration(125000),
				defaultTrack.Info.Length.String(),
				""),
		},
		{
			name:     "with_one_second_remaining_returns_appropriate_message",
			loopMode: enum.LoopOff,
			track:    defaultTrack,
			position: lavalink.Duration(3599999), // 59 min, 59 sec
			expectedMessage: fmt.Sprintf("Now Playing:\n\t*%s* by **%s**\n\t%s `[%s / %s]`\n\t%s",
				defaultTrack.Info.Title,
				defaultTrack.Info.Author,
				"1 sec remaining.",
				lavalink.Duration(3599999),
				defaultTrack.Info.Length.String(),
				""),
		},
		{
			name:     "with_loop_single_returns_appropriate_message",
			loopMode: enum.LoopTrack,
			track:    defaultTrack,
			position: lavalink.Minute,
			expectedMessage: fmt.Sprintf("Now Playing:\n\t*%s* by **%s**\n\t%s `[%s / %s]`\n\t%s",
				defaultTrack.Info.Title,
				defaultTrack.Info.Author,
				"59 mins, 0 secs remaining.",
				lavalink.Minute,
				defaultTrack.Info.Length.String(),
				"**Looping Current Track**\n"),
		},
		{
			name:     "with_loop_queue_returns_appropriate_message",
			loopMode: enum.LoopQueue,
			track:    defaultTrack,
			position: lavalink.Minute,
			expectedMessage: fmt.Sprintf("Now Playing:\n\t*%s* by **%s**\n\t%s `[%s / %s]`\n\t%s",
				defaultTrack.Info.Title,
				defaultTrack.Info.Author,
				"59 mins, 0 secs remaining.",
				lavalink.Minute,
				defaultTrack.Info.Length.String(),
				"**Looping Queue**\n"),
		},
	}

	cmd := command.NewNowPlayingCmd()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc := tc
			t.Parallel()

			fakeMusicPlayer := NewDefaultFakeMusicPlayer()
			fakeMessenger := NewFakeMessenger()

			fakeMusicPlayer.ErrorsToReturn = tc.errorsFromPlayer
			fakeMusicPlayer.loopMode = tc.loopMode
			fakeMusicPlayer.LoadedTrack = &tc.track
			fakeMusicPlayer.CurrentPosition = tc.position

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
