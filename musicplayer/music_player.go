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
	GuildID         snowflake.ID
	Player          lavalink.Player
	Queue           MusicQueue
	SearchResults   SearchResults
	LoopMode        LoopMode
	GotDisconnected bool
}

func NewMusicPlayer(guildId snowflake.ID, link disgolink.Link) *MusicPlayer {
	musicPlayer := MusicPlayer{
		GuildID:         guildId,
		Queue:           NewMusicQueue(),
		SearchResults:   NewSearchResults(),
		LoopMode:        LoopOff,
		GotDisconnected: false,
	}

	musicPlayer.CreatePlayer(link)
	return &musicPlayer
}

func (mp *MusicPlayer) CreatePlayer(link disgolink.Link) {
	player := link.Player(mp.GuildID)
	player.SetVolume(42)
	player.AddListener(NewEventListener(&mp.Queue, &mp.LoopMode, &mp.GotDisconnected))
	mp.Player = player
}

func (mp *MusicPlayer) RecreatePlayer(link disgolink.Link) error {
	err := mp.Player.Destroy()
	if err != nil {
		return err
	}

	mp.CreatePlayer(link)
	return nil
}
