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
	Lavalink        *disgolink.Link
	Player          lavalink.Player
	Queue           MusicQueue
	searchResults   SearchResults
	LoopMode        LoopMode
	GotDisconnected bool
}

func NewMusicPlayer(guildId snowflake.ID, link disgolink.Link) *MusicPlayer {
	musicPlayer := MusicPlayer{
		GuildID:         guildId,
		Player:          link.Player(guildId),
		Queue:           NewMusicQueue(),
		searchResults:   NewSearchResults(),
		LoopMode:        LoopOff,
		GotDisconnected: false,
	}

	musicPlayer.Player.SetVolume(42)
	musicPlayer.Player.AddListener(NewEventListener(&musicPlayer.Queue, &musicPlayer.LoopMode, &musicPlayer.GotDisconnected))

	return &musicPlayer
}

func (mp *MusicPlayer) RecreatePlayer(link disgolink.Link) error {
	err := mp.Player.Destroy()
	if err != nil {
		return err
	}

	player := link.Player(mp.GuildID)
	player.SetVolume(42)
	player.AddListener(NewEventListener(&mp.Queue, &mp.LoopMode, &mp.GotDisconnected))
	mp.Player = player
	return nil
}

func (mp *MusicPlayer) SearchResults() *SearchResults {
	return &mp.searchResults
}
