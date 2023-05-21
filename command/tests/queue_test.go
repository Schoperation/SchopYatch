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

	//defaultQueue := createFakeQueue(2)
	//longQueue := createFakeQueue(15)

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
