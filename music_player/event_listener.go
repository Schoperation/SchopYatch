package music_player

import (
	"context"
	"log"

	"github.com/disgoorg/disgolink/v2/disgolink"
	"github.com/disgoorg/disgolink/v2/lavalink"
	"github.com/disgoorg/snowflake/v2"
)

type MusicPlayerEventListener struct {
	musicPlayers *map[snowflake.ID]*MusicPlayer
}

func NewMusicPlayerEventListener(musicPlayers *map[snowflake.ID]*MusicPlayer) MusicPlayerEventListener {
	return MusicPlayerEventListener{
		musicPlayers: musicPlayers,
	}
}

func (listener *MusicPlayerEventListener) OnTrackEnd(player disgolink.Player, event lavalink.TrackEndEvent) {
	musicPlayerMap := *listener.musicPlayers
	musicPlayer := musicPlayerMap[player.GuildID()]

	if !event.Reason.MayStartNext() || (musicPlayer.IsQueueEmpty() && !musicPlayer.IsLoopModeTrack()) {
		return
	}

	finishedTrack, err := player.Lavalink().BestNode().DecodeTrack(context.TODO(), event.EncodedTrack)
	if err != nil {
		log.Printf("Error occurred decoding the finished track: %v\n", err)
		return
	}

	var nextTrack *lavalink.Track
	if musicPlayer.IsLoopModeTrack() {
		finishedTrackCopy := *finishedTrack
		nextTrack = &finishedTrackCopy
	} else {
		nextTrack, err = musicPlayer.RemoveNextTrackFromQueue()
		if err != nil {
			log.Printf("Error retrieving next track from queue: %v\n", err)
			return
		}
	}

	if nextTrack == nil {
		log.Printf("The next track is nil!\n")
		return
	}

	_, err = musicPlayer.Load(*nextTrack)
	if err != nil {
		log.Printf("Error occurred loading nextTrack: %v\n", err)
		return
	}

	if musicPlayer.IsLoopModeQueue() {
		_, err := musicPlayer.Load(*finishedTrack)
		if err != nil {
			log.Printf("Error occurred re-queueing finishedTrack: %v\n", err)
		}
	}
}

func (listener *MusicPlayerEventListener) OnWebSocketClosed(player disgolink.Player, event lavalink.WebSocketClosedEvent) {
	musicPlayerMap := *listener.musicPlayers
	musicPlayer := musicPlayerMap[player.GuildID()]

	//log.Printf("Socket closed: %s", event.Reason)

	musicPlayer.disconnected = true
}
