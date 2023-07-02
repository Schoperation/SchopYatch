package command_tests

import (
	"errors"
	"fmt"
	"schoperation/schopyatch/command"
	"schoperation/schopyatch/enum"
	"schoperation/schopyatch/msg"
	"schoperation/schopyatch/music_player"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestQueueCmd(t *testing.T) {
	defaultTrack := music_player.NewTrack("test", "title", "author", time.Hour)

	shortQueue := createFakeQueue(1)
	defaultQueue := createFakeQueue(2)
	longQueue := createFakeQueue(15)

	testCases := []struct {
		name             string
		inputOpts        []string
		playerConfig     playerConfig
		errorsFromPlayer map[string]error
		expectedMessage  string
	}{
		{
			name: "empty_queue_with_no_loaded_track_returns_appropriate_success_message",
			playerConfig: playerConfig{
				queue: music_player.NewMusicQueue(),
			},
			errorsFromPlayer: map[string]error{"GetLoadedTrack": errors.New(msg.NoLoadedTrack)},
			expectedMessage:  "Queue is empty.\n",
		},
		{
			name: "empty_queue_with_loaded_track_returns_appropriate_success_message",
			playerConfig: playerConfig{
				queue:           music_player.NewMusicQueue(),
				loadedTrack:     &defaultTrack,
				currentPosition: time.Minute,
			},
			expectedMessage: fmt.Sprintf("Now Playing:\n\t*%s* by **%s** `[%s / %s]`\n\n", defaultTrack.Title(), defaultTrack.Author(), time.Minute, defaultTrack.Length().String()) +
				"Queue is empty.\n",
		},
		{
			name: "one_sized_queue_returns_appropriate_success_message",
			playerConfig: playerConfig{
				queue:           shortQueue,
				loadedTrack:     &defaultTrack,
				currentPosition: time.Minute,
			},
			expectedMessage: fmt.Sprintf("Now Playing:\n\t*%s* by **%s** `[%s / %s]`\n\n", defaultTrack.Title(), defaultTrack.Author(), time.Minute, defaultTrack.Length().String()) +
				fmt.Sprintf("Total of **%d** track in the queue. `[%s]`\n", 1, shortQueue.Duration().String()) +
				fmt.Sprintf("Page **%d** of **%d**:\n\n", 1, 1) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 1, shortQueue.PeekAt(0).Title(), shortQueue.PeekAt(0).Author(), shortQueue.PeekAt(0).Length()),
		},
		{
			name: "two_sized_queue_returns_appropriate_success_message",
			playerConfig: playerConfig{
				queue:           defaultQueue,
				loadedTrack:     &defaultTrack,
				currentPosition: time.Minute,
			},
			expectedMessage: fmt.Sprintf("Now Playing:\n\t*%s* by **%s** `[%s / %s]`\n\n", defaultTrack.Title(), defaultTrack.Author(), time.Minute, defaultTrack.Length().String()) +
				fmt.Sprintf("Total of **%d** tracks in the queue. `[%s]`\n", 2, defaultQueue.Duration().String()) +
				fmt.Sprintf("Page **%d** of **%d**:\n\n", 1, 1) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 1, defaultQueue.PeekAt(0).Title(), defaultQueue.PeekAt(0).Author(), defaultQueue.PeekAt(0).Length()) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 2, defaultQueue.PeekAt(1).Title(), defaultQueue.PeekAt(1).Author(), defaultQueue.PeekAt(1).Length()),
		},
		{
			name: "multi_sized_queue_returns_appropriate_success_message",
			playerConfig: playerConfig{
				queue:           longQueue,
				loadedTrack:     &defaultTrack,
				currentPosition: time.Minute,
			},
			expectedMessage: fmt.Sprintf("Now Playing:\n\t*%s* by **%s** `[%s / %s]`\n\n", defaultTrack.Title(), defaultTrack.Author(), time.Minute, defaultTrack.Length().String()) +
				fmt.Sprintf("Total of **%d** tracks in the queue. `[%s]`\n", 15, longQueue.Duration().String()) +
				fmt.Sprintf("Page **%d** of **%d**:\n\n", 1, 2) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 1, longQueue.PeekAt(0).Title(), longQueue.PeekAt(0).Author(), longQueue.PeekAt(0).Length()) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 2, longQueue.PeekAt(1).Title(), longQueue.PeekAt(1).Author(), longQueue.PeekAt(1).Length()) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 3, longQueue.PeekAt(2).Title(), longQueue.PeekAt(2).Author(), longQueue.PeekAt(2).Length()) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 4, longQueue.PeekAt(3).Title(), longQueue.PeekAt(3).Author(), longQueue.PeekAt(3).Length()) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 5, longQueue.PeekAt(4).Title(), longQueue.PeekAt(4).Author(), longQueue.PeekAt(4).Length()) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 6, longQueue.PeekAt(5).Title(), longQueue.PeekAt(5).Author(), longQueue.PeekAt(5).Length()) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 7, longQueue.PeekAt(6).Title(), longQueue.PeekAt(6).Author(), longQueue.PeekAt(6).Length()) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 8, longQueue.PeekAt(7).Title(), longQueue.PeekAt(7).Author(), longQueue.PeekAt(7).Length()) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 9, longQueue.PeekAt(8).Title(), longQueue.PeekAt(8).Author(), longQueue.PeekAt(8).Length()) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 10, longQueue.PeekAt(9).Title(), longQueue.PeekAt(9).Author(), longQueue.PeekAt(9).Length()),
		},
		{
			name:      "multi_sized_queue_with_page_two_selected_returns_appropriate_success_message",
			inputOpts: []string{"2"},
			playerConfig: playerConfig{
				queue:           longQueue,
				loadedTrack:     &defaultTrack,
				currentPosition: time.Minute,
			},
			expectedMessage: fmt.Sprintf("Now Playing:\n\t*%s* by **%s** `[%s / %s]`\n\n", defaultTrack.Title(), defaultTrack.Author(), time.Minute, defaultTrack.Length().String()) +
				fmt.Sprintf("Total of **%d** tracks in the queue. `[%s]`\n", 15, longQueue.Duration().String()) +
				fmt.Sprintf("Page **%d** of **%d**:\n\n", 2, 2) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 11, longQueue.PeekAt(10).Title(), longQueue.PeekAt(10).Author(), longQueue.PeekAt(10).Length()) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 12, longQueue.PeekAt(11).Title(), longQueue.PeekAt(11).Author(), longQueue.PeekAt(11).Length()) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 13, longQueue.PeekAt(12).Title(), longQueue.PeekAt(12).Author(), longQueue.PeekAt(12).Length()) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 14, longQueue.PeekAt(13).Title(), longQueue.PeekAt(13).Author(), longQueue.PeekAt(13).Length()) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 15, longQueue.PeekAt(14).Title(), longQueue.PeekAt(14).Author(), longQueue.PeekAt(14).Length()),
		},
		{
			name:      "multi_sized_queue_with_out_of_bounds_page_selected_returns_appropriate_success_message",
			inputOpts: []string{"99999"},
			playerConfig: playerConfig{
				queue:           longQueue,
				loadedTrack:     &defaultTrack,
				currentPosition: time.Minute,
			},
			expectedMessage: fmt.Sprintf("Now Playing:\n\t*%s* by **%s** `[%s / %s]`\n\n", defaultTrack.Title(), defaultTrack.Author(), time.Minute, defaultTrack.Length().String()) +
				fmt.Sprintf("Total of **%d** tracks in the queue. `[%s]`\n", 15, longQueue.Duration().String()) +
				fmt.Sprintf("Page **%d** of **%d**:\n\n", 1, 2) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 1, longQueue.PeekAt(0).Title(), longQueue.PeekAt(0).Author(), longQueue.PeekAt(0).Length()) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 2, longQueue.PeekAt(1).Title(), longQueue.PeekAt(1).Author(), longQueue.PeekAt(1).Length()) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 3, longQueue.PeekAt(2).Title(), longQueue.PeekAt(2).Author(), longQueue.PeekAt(2).Length()) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 4, longQueue.PeekAt(3).Title(), longQueue.PeekAt(3).Author(), longQueue.PeekAt(3).Length()) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 5, longQueue.PeekAt(4).Title(), longQueue.PeekAt(4).Author(), longQueue.PeekAt(4).Length()) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 6, longQueue.PeekAt(5).Title(), longQueue.PeekAt(5).Author(), longQueue.PeekAt(5).Length()) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 7, longQueue.PeekAt(6).Title(), longQueue.PeekAt(6).Author(), longQueue.PeekAt(6).Length()) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 8, longQueue.PeekAt(7).Title(), longQueue.PeekAt(7).Author(), longQueue.PeekAt(7).Length()) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 9, longQueue.PeekAt(8).Title(), longQueue.PeekAt(8).Author(), longQueue.PeekAt(8).Length()) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 10, longQueue.PeekAt(9).Title(), longQueue.PeekAt(9).Author(), longQueue.PeekAt(9).Length()),
		},
		{
			name:      "two_sized_queue_with_page_one_selected_returns_appropriate_success_message",
			inputOpts: []string{"1"},
			playerConfig: playerConfig{
				queue:           defaultQueue,
				loadedTrack:     &defaultTrack,
				currentPosition: time.Minute,
			},
			expectedMessage: fmt.Sprintf("Now Playing:\n\t*%s* by **%s** `[%s / %s]`\n\n", defaultTrack.Title(), defaultTrack.Author(), time.Minute, defaultTrack.Length().String()) +
				fmt.Sprintf("Total of **%d** tracks in the queue. `[%s]`\n", 2, defaultQueue.Duration().String()) +
				fmt.Sprintf("Page **%d** of **%d**:\n\n", 1, 1) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 1, defaultQueue.PeekAt(0).Title(), defaultQueue.PeekAt(0).Author(), defaultQueue.PeekAt(0).Length()) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 2, defaultQueue.PeekAt(1).Title(), defaultQueue.PeekAt(1).Author(), defaultQueue.PeekAt(1).Length()),
		},
		{
			name: "two_sized_queue_with_loop_mode_track_returns_appropriate_success_message",
			playerConfig: playerConfig{
				queue:           defaultQueue,
				loadedTrack:     &defaultTrack,
				currentPosition: time.Minute,
				loopMode:        enum.LoopTrack,
			},
			expectedMessage: fmt.Sprintf("Now Playing:\n\t*%s* by **%s** `[%s / %s]`\n\n", defaultTrack.Title(), defaultTrack.Author(), time.Minute, defaultTrack.Length().String()) +
				"**Looping Current Track**\n" +
				fmt.Sprintf("Total of **%d** tracks in the queue. `[%s]`\n", 2, defaultQueue.Duration().String()) +
				fmt.Sprintf("Page **%d** of **%d**:\n\n", 1, 1) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 1, defaultQueue.PeekAt(0).Title(), defaultQueue.PeekAt(0).Author(), defaultQueue.PeekAt(0).Length()) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 2, defaultQueue.PeekAt(1).Title(), defaultQueue.PeekAt(1).Author(), defaultQueue.PeekAt(1).Length()),
		},
		{
			name: "two_sized_queue_with_loop_mode_queue_returns_appropriate_success_message",
			playerConfig: playerConfig{
				queue:           defaultQueue,
				loadedTrack:     &defaultTrack,
				currentPosition: time.Minute,
				loopMode:        enum.LoopQueue,
			},
			expectedMessage: fmt.Sprintf("Now Playing:\n\t*%s* by **%s** `[%s / %s]`\n\n", defaultTrack.Title(), defaultTrack.Author(), time.Minute, defaultTrack.Length().String()) +
				"**Looping Queue**\n" +
				fmt.Sprintf("Total of **%d** tracks in the queue. `[%s]`\n", 2, defaultQueue.Duration().String()) +
				fmt.Sprintf("Page **%d** of **%d**:\n\n", 1, 1) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 1, defaultQueue.PeekAt(0).Title(), defaultQueue.PeekAt(0).Author(), defaultQueue.PeekAt(0).Length()) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 2, defaultQueue.PeekAt(1).Title(), defaultQueue.PeekAt(1).Author(), defaultQueue.PeekAt(1).Length()),
		},
	}

	cmd := command.NewQueueCmd()

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
