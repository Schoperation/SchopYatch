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
	searchResults   []lavalink.AudioTrack
	LoopMode        LoopMode
	GotDisconnected bool
}

func NewMusicPlayer(guildId snowflake.ID, link disgolink.Link) *MusicPlayer {
	musicPlayer := MusicPlayer{
		GuildID:         guildId,
		Player:          link.Player(guildId),
		Queue:           NewMusicQueue(),
		searchResults:   []lavalink.AudioTrack{},
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

func (mp *MusicPlayer) AddSearchResult(track lavalink.AudioTrack) {
	mp.searchResults = append(mp.searchResults, track)
}

func (mp *MusicPlayer) GetSearchResult(i int) *lavalink.AudioTrack {
	if i < 0 || i >= len(mp.searchResults) {
		return nil
	}

	return &mp.searchResults[i]
}

func (mp *MusicPlayer) GetLengthOfSearchResults() int {
	return len(mp.searchResults)
}

func (mp *MusicPlayer) ClearSearchResults() {
	mp.searchResults = []lavalink.AudioTrack{}
}
