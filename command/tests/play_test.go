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

type playerConfig struct {
	isPlayerPaused bool
	loadedTrack    lavalink.Track
	searchResults  music_player.SearchResults
	queue          music_player.MusicQueue
}

func TestPlayCmd(t *testing.T) {
	defaultTrack := lavalink.Track{
		Encoded: "test",
		Info: lavalink.TrackInfo{
			Length: lavalink.Hour,
			Author: "author",
			Title:  "title",
		},
	}

	defaultSearchResults := music_player.NewSearchResults()
	defaultSearchResults.AddResults([]lavalink.Track{
		defaultTrack,
		{
			Encoded: "test2",
			Info: lavalink.TrackInfo{
				Length: lavalink.Hour,
				Author: "author2",
				Title:  "title2",
			},
		},
		{
			Encoded: "test3",
			Info: lavalink.TrackInfo{
				Length: lavalink.Hour,
				Author: "author3",
				Title:  "title3",
			},
		},
	})

	testCases := []struct {
		name             string
		inputOpts        []string
		playerConfig     playerConfig
		statusFromPlayer enum.PlayerStatus
		errorFromPlayer  error
		expectedMessage  string
	}{
		{
			name: "with_paused_player_unpauses_successfully",
			playerConfig: playerConfig{
				isPlayerPaused: true,
			},
			statusFromPlayer: enum.StatusSuccess,
			expectedMessage:  "Resuming.",
		},
		{
			name: "with_unpaused_player_returns_appropriate_error_message",
			playerConfig: playerConfig{
				isPlayerPaused: false,
			},
			statusFromPlayer: enum.StatusAlreadyUnpaused,
			expectedMessage:  "Bruh where's your song??",
		},
		{
			name:      "with_no_search_results_returns_appropriate_error_message",
			inputOpts: []string{"1"},
			playerConfig: playerConfig{
				searchResults: music_player.NewSearchResults(),
			},
			expectedMessage: fmt.Sprintf("Selected thin air. Try a number between 1 and %d.", 0),
		},
		{
			name:      "with_out_of_bounds_search_result_returns_appropriate_error_message",
			inputOpts: []string{"4"},
			playerConfig: playerConfig{
				searchResults: defaultSearchResults,
			},
			expectedMessage: fmt.Sprintf("Selected thin air. Try a number between 1 and %d.", 3),
		},
		{
			name:      "with_search_result_but_no_voice_state_returns_appropriate_error_message",
			inputOpts: []string{"2"},
			playerConfig: playerConfig{
				searchResults: defaultSearchResults,
			},
			errorFromPlayer: errors.New(msg.VoiceStateNotFound),
			expectedMessage: "Dude you're not in a voice channel... get in one I can see!",
		},
		{
			name:      "with_search_result_and_no_loaded_track_returns_appropriate_success_message",
			inputOpts: []string{"2"},
			playerConfig: playerConfig{
				searchResults: defaultSearchResults,
			},
			statusFromPlayer: enum.StatusSuccess,
			expectedMessage:  fmt.Sprintf("Now playing *%s* by **%s**.", "title2", "author2"),
		},
		{
			name:      "with_search_result_and_loaded_track_returns_appropriate_success_message",
			inputOpts: []string{"2"},
			playerConfig: playerConfig{
				searchResults: defaultSearchResults,
				loadedTrack:   defaultTrack,
			},
			statusFromPlayer: enum.StatusQueued,
			expectedMessage:  fmt.Sprintf("Queued *%s* by **%s**.", "title2", "author2"),
		},
	}

	cmd := command.NewPlayCmd()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc := tc
			t.Parallel()

			fakeMusicPlayer := NewDefaultFakeMusicPlayer()
			fakeMessenger := NewFakeMessenger()

			fakeMusicPlayer.ErrorToReturn = tc.errorFromPlayer
			fakeMusicPlayer.StatusToReturn = tc.statusFromPlayer

			fakeMusicPlayer.Paused = tc.playerConfig.isPlayerPaused
			fakeMusicPlayer.LoadedTrack = &tc.playerConfig.loadedTrack
			fakeMusicPlayer.searchResults = tc.playerConfig.searchResults
			fakeMusicPlayer.queue = tc.playerConfig.queue

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
