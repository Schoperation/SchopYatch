package command_tests

import (
	"errors"
	"fmt"
	"schoperation/schopyatch/command"
	"schoperation/schopyatch/msg"
	"schoperation/schopyatch/music_player"
	"testing"
	"time"

	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/stretchr/testify/require"
)

func TestSeekCmd(t *testing.T) {
	defaultTrack := music_player.NewTrack("test", "title", "author", time.Hour*2)

	testCases := []struct {
		name             string
		inputOpts        []string
		playerConfig     playerConfig
		errorsFromPlayer map[string]error
		expectedMessage  string
	}{
		{
			name:            "no_selection_returns_appropriate_error_message",
			expectedMessage: fmt.Sprintf("Need to specify a time to seek to. `%sseek hh:mm:ss`", ";;"),
		},
		{
			name:            "invalid_selection_returns_appropriate_error_message",
			inputOpts:       []string{"ds:f2:3f"},
			expectedMessage: fmt.Sprintf("Detected something that ain't a number... cmon man... `%sseek hh:mm:ss`", ";;"),
		},
		{
			name:             "no_loaded_track_returns_appropriate_error_message",
			inputOpts:        []string{"30"},
			errorsFromPlayer: map[string]error{"Seek": errors.New(msg.NoLoadedTrack)},
			expectedMessage:  "Can't seek through nothing... not a sikh move if you ask me...",
		},
		{
			name:      "negative_duration_returns_appropriate_error_message",
			inputOpts: []string{"-30"},
			playerConfig: playerConfig{
				loadedTrack: &defaultTrack,
			},
			errorsFromPlayer: map[string]error{"Seek": errors.New(msg.NegativeDuration)},
			expectedMessage:  "Can't seek back in time. Try adopting a more positive outlook.",
		},
		{
			name:      "out_of_bounds_duration_returns_appropriate_error_message",
			inputOpts: []string{"99:99:99"},
			playerConfig: playerConfig{
				loadedTrack: &defaultTrack,
			},
			errorsFromPlayer: map[string]error{"Seek": errors.New(msg.IndexOutOfBounds)},
			expectedMessage:  fmt.Sprintf("Can't seek beyond the track length. It's %s long.", defaultTrack.Length().String()),
		},
		{
			name:      "seconds_duration_returns_appropriate_success_message",
			inputOpts: []string{"30"},
			playerConfig: playerConfig{
				loadedTrack: &defaultTrack,
			},
			expectedMessage: fmt.Sprintf("Seeked to `%s`", lavalink.Duration(30000).String()),
		},
		{
			name:      "minutes_duration_returns_appropriate_success_message",
			inputOpts: []string{"1:30"},
			playerConfig: playerConfig{
				loadedTrack: &defaultTrack,
			},
			expectedMessage: fmt.Sprintf("Seeked to `%s`", lavalink.Duration(90000).String()),
		},
		{
			name:      "hours_duration_returns_appropriate_success_message",
			inputOpts: []string{"01:01:30"},
			playerConfig: playerConfig{
				loadedTrack: &defaultTrack,
			},
			expectedMessage: fmt.Sprintf("Seeked to `%s`", lavalink.Duration(3690000).String()),
		},
	}

	cmd := command.NewSeekCmd()

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
		})
	}
}
