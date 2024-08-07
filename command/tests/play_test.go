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

func TestPlayCmd(t *testing.T) {
	defaultTrack := music_player.NewTrack("test", "title", "author", time.Hour)

	defaultSearchResults := music_player.NewSearchResults()
	defaultSearchResults.AddResults([]music_player.Track{
		defaultTrack,
		music_player.NewTrack("test2", "title2", "author2", time.Hour),
		music_player.NewTrack("test3", "title3", "author3", time.Hour),
	})

	testCases := []struct {
		name             string
		inputOpts        []string
		playerConfig     playerConfig
		statusFromPlayer enum.PlayerStatus
		errorsFromPlayer map[string]error
		expectedMessage  string
	}{
		{
			name: "with_paused_player_unpauses_successfully",
			playerConfig: playerConfig{
				isPlayerPaused: true,
				loadedTrack:    &defaultTrack,
			},
			statusFromPlayer: enum.StatusSuccess,
			expectedMessage:  "Resuming.",
		},
		{
			name: "with_unpaused_player_returns_appropriate_error_message",
			playerConfig: playerConfig{
				isPlayerPaused: false,
				loadedTrack:    &defaultTrack,
			},
			statusFromPlayer: enum.StatusAlreadyUnpaused,
			expectedMessage:  "Bruh where's your song??",
		},
		{
			name: "with_no_track_loaded_returns_appropriate_error_message",
			playerConfig: playerConfig{
				isPlayerPaused: false,
				loadedTrack:    nil,
			},
			statusFromPlayer: enum.StatusFailed,
			expectedMessage:  "Bruh where's your song??",
		},
		{
			name:      "with_search_selection_but_no_search_results_returns_appropriate_error_message",
			inputOpts: []string{"1"},
			playerConfig: playerConfig{
				searchResults: music_player.NewSearchResults(),
			},
			expectedMessage: fmt.Sprintf("Selected thin air. Try a number between 1 and %d.", 0),
		},
		{
			name:      "with_search_selection_but_out_of_bounds_returns_appropriate_error_message",
			inputOpts: []string{"4"},
			playerConfig: playerConfig{
				searchResults: defaultSearchResults,
			},
			expectedMessage: fmt.Sprintf("Selected thin air. Try a number between 1 and %d.", 3),
		},
		{
			name:      "with_search_selection_but_no_voice_state_returns_appropriate_error_message",
			inputOpts: []string{"2"},
			playerConfig: playerConfig{
				searchResults: defaultSearchResults,
			},
			errorsFromPlayer: map[string]error{"JoinVoiceChannel": errors.New(msg.VoiceStateNotFound)},
			expectedMessage:  "Dude you're not in a voice channel... get in one I can see!",
		},
		{
			name:      "with_search_selection_and_no_loaded_track_returns_appropriate_success_message",
			inputOpts: []string{"2"},
			playerConfig: playerConfig{
				searchResults: defaultSearchResults,
			},
			statusFromPlayer: enum.StatusSuccess,
			expectedMessage:  fmt.Sprintf("Now playing *%s* by **%s**.", "title2", "author2"),
		},
		{
			name:      "with_search_selection_and_loaded_track_returns_appropriate_success_message",
			inputOpts: []string{"2"},
			playerConfig: playerConfig{
				searchResults: defaultSearchResults,
				loadedTrack:   &defaultTrack,
			},
			statusFromPlayer: enum.StatusQueued,
			expectedMessage:  fmt.Sprintf("Queued *%s* by **%s**.", "title2", "author2"),
		},
		{
			name:      "with_query_but_no_results_found_returns_appropriate_error_message",
			inputOpts: []string{"ace", "attorney", "all", "pursuit", "themes"},
			playerConfig: playerConfig{
				searchResults: defaultSearchResults,
			},
			errorsFromPlayer: map[string]error{"ProcessQuery": errors.New(msg.NoResultsFound)},
			expectedMessage:  "No results. Try some other keywords? Such as OFFICIAL, FEATURING, ft., THE TRUTH ABOUT, IS A FRAUD, or CHARLIE",
		},
		{
			name:      "with_url_track_queued_returns_appropriate_success_message",
			inputOpts: []string{"https://www.youtube.com/watch?v=enuOArEfqGo"},
			playerConfig: playerConfig{
				searchResults: defaultSearchResults,
				loadedTrack:   &defaultTrack,
			},
			statusFromPlayer: enum.StatusQueued,
			expectedMessage:  fmt.Sprintf("Queued *%s* by **%s**.", "title", "author"),
		},
		{
			name:      "with_empty_url_playlist_queued_returns_appropriate_success_message",
			inputOpts: []string{"https://www.youtube.com/watch?v=enuOArEfqGo"},
			playerConfig: playerConfig{
				searchResults: defaultSearchResults,
				loadedTrack:   &defaultTrack,
				tracksQueued:  0,
			},
			statusFromPlayer: enum.StatusQueuedList,
			expectedMessage:  "Queued nothing. What the...?",
		},
		{
			name:      "with_url_playlist_with_one_track_queued_returns_appropriate_success_message",
			inputOpts: []string{"https://www.youtube.com/watch?v=enuOArEfqGo"},
			playerConfig: playerConfig{
				searchResults: defaultSearchResults,
				loadedTrack:   &defaultTrack,
				tracksQueued:  1,
			},
			statusFromPlayer: enum.StatusQueuedList,
			expectedMessage:  "Queued **1** additional track. Just a one hit wonder, huh?",
		},
		{
			name:      "with_url_playlist_with_multiple_tracks_queued_returns_appropriate_success_message",
			inputOpts: []string{"https://www.youtube.com/watch?v=enuOArEfqGo"},
			playerConfig: playerConfig{
				searchResults: defaultSearchResults,
				loadedTrack:   &defaultTrack,
				tracksQueued:  5,
			},
			statusFromPlayer: enum.StatusQueuedList,
			expectedMessage:  fmt.Sprintf("Queued **%d** additional tracks.", 5),
		},
		{
			name:      "with_url_playlist_with_playing_track_and_one_track_queued_returns_appropriate_success_message",
			inputOpts: []string{"https://www.youtube.com/watch?v=enuOArEfqGo"},
			playerConfig: playerConfig{
				searchResults: defaultSearchResults,
				loadedTrack:   &defaultTrack,
				tracksQueued:  1,
			},
			statusFromPlayer: enum.StatusPlayingAndQueuedList,
			expectedMessage:  fmt.Sprintf("Now playing *%s* by **%s**.\nQueued **1** additional track.", "title", "author"),
		},
		{
			name:      "with_query_and_results_found_returns_appropriate_success_message",
			inputOpts: []string{"ace", "attorney", "all", "pursuit", "themes"},
			playerConfig: playerConfig{
				searchResults: defaultSearchResults,
			},
			statusFromPlayer: enum.StatusSearchSuccess,
			expectedMessage: "Search Results:\n\n" +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 1, "title", "author", time.Hour) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 2, "title2", "author2", time.Hour) +
				fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", 3, "title3", "author3", time.Hour) +
				fmt.Sprintf("\nUse `%splay n` to pick a track to play. Results will be available until the next query.", ";;"),
		},
		{
			name:      "with_url_and_no_loaded_track_returns_appropriate_success_message",
			inputOpts: []string{"https://www.youtube.com/watch?v=enuOArEfqGo"},
			playerConfig: playerConfig{
				searchResults: defaultSearchResults,
				loadedTrack:   &defaultTrack,
			},
			statusFromPlayer: enum.StatusSuccess,
			expectedMessage:  fmt.Sprintf("Now playing *%s* by **%s**.", "title", "author"),
		},
		{
			name:      "with_url_and_loaded_track_returns_appropriate_success_message",
			inputOpts: []string{"https://www.youtube.com/watch?v=enuOArEfqGo"},
			playerConfig: playerConfig{
				searchResults: defaultSearchResults,
				loadedTrack:   &defaultTrack,
			},
			statusFromPlayer: enum.StatusQueued,
			expectedMessage:  fmt.Sprintf("Queued *%s* by **%s**.", "title", "author"),
		},
	}

	cmd := command.NewPlayCmd()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc := tc
			t.Parallel()

			fakeMusicPlayer := NewDefaultFakeMusicPlayer()
			fakeMessenger := NewFakeMessenger()

			fakeMusicPlayer.ErrorsToReturn = tc.errorsFromPlayer
			fakeMusicPlayer.StatusToReturn = tc.statusFromPlayer

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
