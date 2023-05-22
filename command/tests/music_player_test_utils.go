package command_tests

import (
	"fmt"
	"schoperation/schopyatch/enum"
	"schoperation/schopyatch/music_player"

	"github.com/disgoorg/disgolink/v2/lavalink"
)

type playerConfig struct {
	isPlayerPaused  bool
	loadedTrack     *lavalink.Track
	currentPosition lavalink.Duration
	searchResults   music_player.SearchResults
	queue           music_player.MusicQueue
	tracksQueued    int
	loopMode        enum.LoopMode
}

func (fmp *fakeMusicPlayer) setPlayerConfig(playerConfig playerConfig) {
	fmp.Paused = playerConfig.isPlayerPaused
	fmp.LoadedTrack = playerConfig.loadedTrack
	fmp.CurrentPosition = playerConfig.currentPosition
	fmp.searchResults = playerConfig.searchResults
	fmp.queue = playerConfig.queue
	fmp.TracksQueued = playerConfig.tracksQueued
	fmp.loopMode = playerConfig.loopMode
}

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
