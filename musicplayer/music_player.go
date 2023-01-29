package musicplayer

import (
	"github.com/disgoorg/disgolink/disgolink"
	"github.com/disgoorg/disgolink/lavalink"
	"github.com/disgoorg/snowflake/v2"
)

type LoopMode int

const (
	LoopOff LoopMode = iota
	LoopTrack
	LoopQueue
)

type MusicPlayer struct {
	GuildID  snowflake.ID
	Player   lavalink.Player
	Queue    MusicQueue
	LoopMode LoopMode
}

func NewMusicPlayer(guildId snowflake.ID, lavalink disgolink.Link) *MusicPlayer {
	player := lavalink.Player(guildId)
	player.SetVolume(42)
	player.AddListener(NewEventListener())

	return &MusicPlayer{
		GuildID:  guildId,
		Player:   lavalink.Player(guildId),
		Queue:    *NewMusicQueue(),
		LoopMode: LoopOff,
	}
}
