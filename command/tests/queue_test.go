package command_tests

import (
	"errors"
	"fmt"
	"schoperation/schopyatch/command"
	"schoperation/schopyatch/enum"
	"schoperation/schopyatch/msg"
	"schoperation/schopyatch/music_player"
	"testing"

	"github.com/disgoorg/disgolink/v2/lavalink"
	"github.com/stretchr/testify/require"
)

func createFakeQueue(numTracks int) music_player.MusicQueue {
	queue := music_player.NewMusicQueue()
	for i := 0; i < numTracks; i++ {
		queue.Enqueue(lavalink.Track{
			Encoded: "test",
			Info: lavalink.TrackInfo{
				Length: lavalink.Hour,
				Author: fmt.Sprintf("author%d", i+1),
				Title:  fmt.Sprintf("title%d", i+1),
			},
		})
	}

	return queue
}

func TestQueueCmd(t *testing.T) {
	defaultTrack := lavalink.Track{
		Encoded: "test",
		Info: lavalink.TrackInfo{
			Length: lavalink.Hour,
			Author: "author",
			Title:  "title",
		},
	}

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
				loadedTrack:     defaultTrack,
				currentPosition: lavalink.Minute,
			},
			expectedMessage: fmt.Sprintf("Now Playing:\n\t*%s* by **%s** `[%s / %s]`\n\n", defaultTrack.Info.Title, defaultTrack.Info.Author, lavalink.Minute, defaultTrack.Info.Length.String()) +
				"Queue is empty.\n",
		},
		{
			name: "one_sized_queue_returns_appropriate_success_message",
			playerConfig: playerConfig{
				queue:           shortQueue,
				loadedTrack:     defaultTrack,
				currentPosition: lavalink.Minute,
			},
			expectedMessage: fmt.Sprintf("Now Playing:\n\t*%s* by **%s** `[%s / %s]`\n\n", defaultTrack.Info.Title, defaultTrack.Info.Author, lavalink.Minute, defaultTrack.Info.Length.String()) +
				fmt.Sprintf("Total of **%d** track in the queue. `[%s]`\n", 1, shortQueue.Duration().String()) +
				fmt.Sprintf("Page **%d** of **%d**:\n\n", 1, 1) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 1, shortQueue.PeekAt(0).Info.Title, shortQueue.PeekAt(0).Info.Author, shortQueue.PeekAt(0).Info.Length),
		},
		{
			name: "two_sized_queue_returns_appropriate_success_message",
			playerConfig: playerConfig{
				queue:           defaultQueue,
				loadedTrack:     defaultTrack,
				currentPosition: lavalink.Minute,
			},
			expectedMessage: fmt.Sprintf("Now Playing:\n\t*%s* by **%s** `[%s / %s]`\n\n", defaultTrack.Info.Title, defaultTrack.Info.Author, lavalink.Minute, defaultTrack.Info.Length.String()) +
				fmt.Sprintf("Total of **%d** tracks in the queue. `[%s]`\n", 2, defaultQueue.Duration().String()) +
				fmt.Sprintf("Page **%d** of **%d**:\n\n", 1, 1) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 1, defaultQueue.PeekAt(0).Info.Title, defaultQueue.PeekAt(0).Info.Author, defaultQueue.PeekAt(0).Info.Length) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 2, defaultQueue.PeekAt(1).Info.Title, defaultQueue.PeekAt(1).Info.Author, defaultQueue.PeekAt(1).Info.Length),
		},
		{
			name: "multi_sized_queue_returns_appropriate_success_message",
			playerConfig: playerConfig{
				queue:           longQueue,
				loadedTrack:     defaultTrack,
				currentPosition: lavalink.Minute,
			},
			expectedMessage: fmt.Sprintf("Now Playing:\n\t*%s* by **%s** `[%s / %s]`\n\n", defaultTrack.Info.Title, defaultTrack.Info.Author, lavalink.Minute, defaultTrack.Info.Length.String()) +
				fmt.Sprintf("Total of **%d** tracks in the queue. `[%s]`\n", 15, longQueue.Duration().String()) +
				fmt.Sprintf("Page **%d** of **%d**:\n\n", 1, 2) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 1, longQueue.PeekAt(0).Info.Title, longQueue.PeekAt(0).Info.Author, longQueue.PeekAt(0).Info.Length) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 2, longQueue.PeekAt(1).Info.Title, longQueue.PeekAt(1).Info.Author, longQueue.PeekAt(1).Info.Length) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 3, longQueue.PeekAt(2).Info.Title, longQueue.PeekAt(2).Info.Author, longQueue.PeekAt(2).Info.Length) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 4, longQueue.PeekAt(3).Info.Title, longQueue.PeekAt(3).Info.Author, longQueue.PeekAt(3).Info.Length) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 5, longQueue.PeekAt(4).Info.Title, longQueue.PeekAt(4).Info.Author, longQueue.PeekAt(4).Info.Length) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 6, longQueue.PeekAt(5).Info.Title, longQueue.PeekAt(5).Info.Author, longQueue.PeekAt(5).Info.Length) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 7, longQueue.PeekAt(6).Info.Title, longQueue.PeekAt(6).Info.Author, longQueue.PeekAt(6).Info.Length) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 8, longQueue.PeekAt(7).Info.Title, longQueue.PeekAt(7).Info.Author, longQueue.PeekAt(7).Info.Length) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 9, longQueue.PeekAt(8).Info.Title, longQueue.PeekAt(8).Info.Author, longQueue.PeekAt(8).Info.Length) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 10, longQueue.PeekAt(9).Info.Title, longQueue.PeekAt(9).Info.Author, longQueue.PeekAt(9).Info.Length),
		},
		{
			name:      "multi_sized_queue_with_page_two_selected_returns_appropriate_success_message",
			inputOpts: []string{"2"},
			playerConfig: playerConfig{
				queue:           longQueue,
				loadedTrack:     defaultTrack,
				currentPosition: lavalink.Minute,
			},
			expectedMessage: fmt.Sprintf("Now Playing:\n\t*%s* by **%s** `[%s / %s]`\n\n", defaultTrack.Info.Title, defaultTrack.Info.Author, lavalink.Minute, defaultTrack.Info.Length.String()) +
				fmt.Sprintf("Total of **%d** tracks in the queue. `[%s]`\n", 15, longQueue.Duration().String()) +
				fmt.Sprintf("Page **%d** of **%d**:\n\n", 2, 2) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 11, longQueue.PeekAt(10).Info.Title, longQueue.PeekAt(10).Info.Author, longQueue.PeekAt(10).Info.Length) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 12, longQueue.PeekAt(11).Info.Title, longQueue.PeekAt(11).Info.Author, longQueue.PeekAt(11).Info.Length) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 13, longQueue.PeekAt(12).Info.Title, longQueue.PeekAt(12).Info.Author, longQueue.PeekAt(12).Info.Length) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 14, longQueue.PeekAt(13).Info.Title, longQueue.PeekAt(13).Info.Author, longQueue.PeekAt(13).Info.Length) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 15, longQueue.PeekAt(14).Info.Title, longQueue.PeekAt(14).Info.Author, longQueue.PeekAt(14).Info.Length),
		},
		{
			name:      "multi_sized_queue_with_out_of_bounds_page_selected_returns_appropriate_success_message",
			inputOpts: []string{"99999"},
			playerConfig: playerConfig{
				queue:           longQueue,
				loadedTrack:     defaultTrack,
				currentPosition: lavalink.Minute,
			},
			expectedMessage: fmt.Sprintf("Now Playing:\n\t*%s* by **%s** `[%s / %s]`\n\n", defaultTrack.Info.Title, defaultTrack.Info.Author, lavalink.Minute, defaultTrack.Info.Length.String()) +
				fmt.Sprintf("Total of **%d** tracks in the queue. `[%s]`\n", 15, longQueue.Duration().String()) +
				fmt.Sprintf("Page **%d** of **%d**:\n\n", 1, 2) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 1, longQueue.PeekAt(0).Info.Title, longQueue.PeekAt(0).Info.Author, longQueue.PeekAt(0).Info.Length) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 2, longQueue.PeekAt(1).Info.Title, longQueue.PeekAt(1).Info.Author, longQueue.PeekAt(1).Info.Length) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 3, longQueue.PeekAt(2).Info.Title, longQueue.PeekAt(2).Info.Author, longQueue.PeekAt(2).Info.Length) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 4, longQueue.PeekAt(3).Info.Title, longQueue.PeekAt(3).Info.Author, longQueue.PeekAt(3).Info.Length) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 5, longQueue.PeekAt(4).Info.Title, longQueue.PeekAt(4).Info.Author, longQueue.PeekAt(4).Info.Length) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 6, longQueue.PeekAt(5).Info.Title, longQueue.PeekAt(5).Info.Author, longQueue.PeekAt(5).Info.Length) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 7, longQueue.PeekAt(6).Info.Title, longQueue.PeekAt(6).Info.Author, longQueue.PeekAt(6).Info.Length) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 8, longQueue.PeekAt(7).Info.Title, longQueue.PeekAt(7).Info.Author, longQueue.PeekAt(7).Info.Length) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 9, longQueue.PeekAt(8).Info.Title, longQueue.PeekAt(8).Info.Author, longQueue.PeekAt(8).Info.Length) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 10, longQueue.PeekAt(9).Info.Title, longQueue.PeekAt(9).Info.Author, longQueue.PeekAt(9).Info.Length),
		},
		{
			name:      "two_sized_queue_with_page_one_selected_returns_appropriate_success_message",
			inputOpts: []string{"1"},
			playerConfig: playerConfig{
				queue:           defaultQueue,
				loadedTrack:     defaultTrack,
				currentPosition: lavalink.Minute,
			},
			expectedMessage: fmt.Sprintf("Now Playing:\n\t*%s* by **%s** `[%s / %s]`\n\n", defaultTrack.Info.Title, defaultTrack.Info.Author, lavalink.Minute, defaultTrack.Info.Length.String()) +
				fmt.Sprintf("Total of **%d** tracks in the queue. `[%s]`\n", 2, defaultQueue.Duration().String()) +
				fmt.Sprintf("Page **%d** of **%d**:\n\n", 1, 1) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 1, defaultQueue.PeekAt(0).Info.Title, defaultQueue.PeekAt(0).Info.Author, defaultQueue.PeekAt(0).Info.Length) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 2, defaultQueue.PeekAt(1).Info.Title, defaultQueue.PeekAt(1).Info.Author, defaultQueue.PeekAt(1).Info.Length),
		},
		{
			name: "two_sized_queue_with_loop_mode_track_returns_appropriate_success_message",
			playerConfig: playerConfig{
				queue:           defaultQueue,
				loadedTrack:     defaultTrack,
				currentPosition: lavalink.Minute,
				loopMode:        enum.LoopTrack,
			},
			expectedMessage: fmt.Sprintf("Now Playing:\n\t*%s* by **%s** `[%s / %s]`\n\n", defaultTrack.Info.Title, defaultTrack.Info.Author, lavalink.Minute, defaultTrack.Info.Length.String()) +
				"**Looping Current Track**\n" +
				fmt.Sprintf("Total of **%d** tracks in the queue. `[%s]`\n", 2, defaultQueue.Duration().String()) +
				fmt.Sprintf("Page **%d** of **%d**:\n\n", 1, 1) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 1, defaultQueue.PeekAt(0).Info.Title, defaultQueue.PeekAt(0).Info.Author, defaultQueue.PeekAt(0).Info.Length) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 2, defaultQueue.PeekAt(1).Info.Title, defaultQueue.PeekAt(1).Info.Author, defaultQueue.PeekAt(1).Info.Length),
		},
		{
			name: "two_sized_queue_with_loop_mode_queue_returns_appropriate_success_message",
			playerConfig: playerConfig{
				queue:           defaultQueue,
				loadedTrack:     defaultTrack,
				currentPosition: lavalink.Minute,
				loopMode:        enum.LoopQueue,
			},
			expectedMessage: fmt.Sprintf("Now Playing:\n\t*%s* by **%s** `[%s / %s]`\n\n", defaultTrack.Info.Title, defaultTrack.Info.Author, lavalink.Minute, defaultTrack.Info.Length.String()) +
				"**Looping Queue**\n" +
				fmt.Sprintf("Total of **%d** tracks in the queue. `[%s]`\n", 2, defaultQueue.Duration().String()) +
				fmt.Sprintf("Page **%d** of **%d**:\n\n", 1, 1) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 1, defaultQueue.PeekAt(0).Info.Title, defaultQueue.PeekAt(0).Info.Author, defaultQueue.PeekAt(0).Info.Length) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 2, defaultQueue.PeekAt(1).Info.Title, defaultQueue.PeekAt(1).Info.Author, defaultQueue.PeekAt(1).Info.Length),
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

			fakeMusicPlayer.Paused = tc.playerConfig.isPlayerPaused
			fakeMusicPlayer.LoadedTrack = &tc.playerConfig.loadedTrack
			fakeMusicPlayer.CurrentPosition = tc.playerConfig.currentPosition
			fakeMusicPlayer.searchResults = tc.playerConfig.searchResults
			fakeMusicPlayer.queue = tc.playerConfig.queue
			fakeMusicPlayer.TracksQueued = tc.playerConfig.tracksQueued
			fakeMusicPlayer.loopMode = tc.playerConfig.loopMode

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
