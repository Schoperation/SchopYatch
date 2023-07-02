package music_player

import (
	"context"
	"log"
	"time"

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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	finishedTrack, err := player.Lavalink().BestNode().DecodeTrack(ctx, event.EncodedTrack)
	if err != nil {
		log.Printf("Error occurred decoding the finished track: %v\n", err)
		return
	}

	var nextTrack *Track
	if musicPlayer.IsLoopModeTrack() {
		finishedTrackCopy := toTrack(*finishedTrack)
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
		_, err := musicPlayer.Load(toTrack(*finishedTrack))
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
