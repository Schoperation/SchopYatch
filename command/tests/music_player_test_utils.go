package command_tests

import (
	"fmt"
	"schoperation/schopyatch/enum"
	"schoperation/schopyatch/music_player"
	"time"
)

type playerConfig struct {
	isPlayerPaused  bool
	loadedTrack     *music_player.Track
	currentPosition time.Duration
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
		queue.Enqueue(music_player.NewTrack("test", fmt.Sprintf("title%d", i+1), fmt.Sprintf("author%d", i+1), time.Hour))
	}

	return queue
}
