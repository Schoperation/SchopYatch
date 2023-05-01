package command_tests

import (
	"schoperation/schopyatch/command"
	"schoperation/schopyatch/enum"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoopCmd(t *testing.T) {
	testCases := []struct {
		name             string
		inputOpts        []string
		currentLoopMode  enum.LoopMode
		expectedLoopMode enum.LoopMode
		expectedMessage  string
	}{
		{
			name:             "no_options_with_no_looping_sets_to_loop_track",
			currentLoopMode:  enum.LoopOff,
			expectedLoopMode: enum.LoopTrack,
			expectedMessage:  "Looping the current track.",
		},
		{
			name:             "no_options_with_looping_track_sets_to_off",
			currentLoopMode:  enum.LoopTrack,
			expectedLoopMode: enum.LoopOff,
			expectedMessage:  "Looping off.",
		},
		{
			name:             "no_options_with_looping_queue_sets_to_off",
			currentLoopMode:  enum.LoopQueue,
			expectedLoopMode: enum.LoopOff,
			expectedMessage:  "Looping off.",
		},
		{
			name:             "single_option_with_no_looping_sets_to_loop_track",
			inputOpts:        []string{"single"},
			currentLoopMode:  enum.LoopOff,
			expectedLoopMode: enum.LoopTrack,
			expectedMessage:  "Looping the current track.",
		},
		{
			name:             "single_option_with_looping_queue_sets_to_loop_track",
			inputOpts:        []string{"single"},
			currentLoopMode:  enum.LoopQueue,
			expectedLoopMode: enum.LoopTrack,
			expectedMessage:  "Looping the current track.",
		},
		{
			name:             "single_option_with_looping_track_sets_to_off",
			inputOpts:        []string{"single"},
			currentLoopMode:  enum.LoopTrack,
			expectedLoopMode: enum.LoopOff,
			expectedMessage:  "Looping off.",
		},
		{
			name:             "queue_option_with_no_looping_sets_to_loop_queue",
			inputOpts:        []string{"queue"},
			currentLoopMode:  enum.LoopOff,
			expectedLoopMode: enum.LoopQueue,
			expectedMessage:  "Looping the whole queue.",
		},
		{
			name:             "queue_option_with_looping_track_sets_to_loop_queue",
			inputOpts:        []string{"queue"},
			currentLoopMode:  enum.LoopTrack,
			expectedLoopMode: enum.LoopQueue,
			expectedMessage:  "Looping the whole queue.",
		},
		{
			name:             "queue_option_with_looping_queue_sets_to_off",
			inputOpts:        []string{"queue"},
			currentLoopMode:  enum.LoopQueue,
			expectedLoopMode: enum.LoopOff,
			expectedMessage:  "Looping off.",
		},
		{
			name:             "all_option_with_no_looping_sets_to_loop_queue",
			inputOpts:        []string{"all"},
			currentLoopMode:  enum.LoopOff,
			expectedLoopMode: enum.LoopQueue,
			expectedMessage:  "Looping the whole queue.",
		},
		{
			name:             "list_option_with_no_looping_sets_to_loop_queue",
			inputOpts:        []string{"list"},
			currentLoopMode:  enum.LoopOff,
			expectedLoopMode: enum.LoopQueue,
			expectedMessage:  "Looping the whole queue.",
		},
		{
			name:             "junk_option_with_no_looping_sets_to_off",
			inputOpts:        []string{"asdfsdafsadfasdf"},
			currentLoopMode:  enum.LoopOff,
			expectedLoopMode: enum.LoopOff,
			expectedMessage:  "Looping off.",
		},
	}

	cmd := command.NewLoopCmd()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc := tc
			t.Parallel()

			fakeMusicPlayer := NewDefaultFakeMusicPlayer()
			fakeMessenger := NewFakeMessenger()

			if tc.currentLoopMode == enum.LoopTrack {
				fakeMusicPlayer.SetLoopModeTrack()
			} else if tc.currentLoopMode == enum.LoopQueue {
				fakeMusicPlayer.SetLoopModeQueue()
			} else {
				fakeMusicPlayer.SetLoopModeOff()
			}

			err := cmd.Execute(command.CommandDependencies{
				MusicPlayer: &fakeMusicPlayer,
				Messenger:   &fakeMessenger,
				Event:       nil,
				Prefix:      ";;",
			}, tc.inputOpts...)

			require.Nil(t, err)
			require.Equal(t, tc.expectedMessage, fakeMessenger.sentMessage)
			require.Equal(t, tc.expectedLoopMode, fakeMusicPlayer.loopMode)
		})
	}
}
