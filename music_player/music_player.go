package music_player

import (
	"context"
	"errors"
	"schoperation/schopyatch/util"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgolink/v2/disgolink"
	"github.com/disgoorg/disgolink/v2/lavalink"
	"github.com/disgoorg/snowflake/v2"
)

type LoopMode int

const (
	LoopOff LoopMode = iota
	LoopTrack
	LoopQueue
)

type MusicPlayer struct {
	guildID        snowflake.ID
	lavalinkClient *disgolink.Client
	player         disgolink.Player
	queue          MusicQueue
	searchResults  SearchResults
	loopMode       LoopMode
	disconnected   bool
}

func NewMusicPlayer(guildId snowflake.ID, lavaLinkClient *disgolink.Client) *MusicPlayer {
	musicPlayer := MusicPlayer{
		guildID:        guildId,
		lavalinkClient: lavaLinkClient,
		queue:          NewMusicQueue(),
		searchResults:  NewSearchResults(),
		loopMode:       LoopOff,
		disconnected:   false,
	}

	musicPlayer.newPlayer()
	return &musicPlayer
}

/////////////////////
// Starters
/////////////////////

func (mp *MusicPlayer) JoinVoiceChannel(botClient *bot.Client, userId snowflake.ID) error {
	voiceState, exists := (*botClient).Caches().VoiceState(mp.guildID, userId)
	if !exists {
		return errors.New(util.VoiceStateNotFound)
	}

	err := (*botClient).UpdateVoiceState(context.TODO(), mp.guildID, voiceState.ChannelID, false, true)
	mp.recreatePlayer()
	return err
}

func (mp *MusicPlayer) LeaveVoiceChannel(botClient *bot.Client, shouldReset bool) error {
	if mp.HasLoadedTrack() {
		err := mp.Stop()
		if err != nil {
			return err
		}
	}

	err := (*botClient).UpdateVoiceState(context.TODO(), mp.guildID, nil, false, false)
	if err != nil {
		return err
	}

	if shouldReset {
		mp.queue.Clear()
		mp.searchResults.Clear()
		mp.SetLoopModeOff()
	}

	mp.disconnected = true
	return nil
}

func (mp *MusicPlayer) ProcessQuery(query string) {

}

/////////////////////
// Basic Player Cmds
/////////////////////

func (mp *MusicPlayer) Load() {

}

func (mp *MusicPlayer) Pause() error {
	if !mp.HasLoadedTrack() {
		return errors.New(util.NoLoadedTrack)
	}

	return mp.player.Update(context.TODO(), lavalink.WithPaused(true))
}

func (mp *MusicPlayer) Unpause() error {
	if !mp.HasLoadedTrack() {
		return errors.New(util.NoLoadedTrack)
	}

	return mp.player.Update(context.TODO(), lavalink.WithPaused(false))
}

func (mp *MusicPlayer) Stop() error {
	if !mp.HasLoadedTrack() {
		return errors.New(util.NoLoadedTrack)
	}

	return mp.player.Update(context.TODO(), lavalink.WithNullTrack())
}

func (mp *MusicPlayer) Seek() {

}

func (mp *MusicPlayer) Skip() {

}

func (mp *MusicPlayer) SkipTo() {

}

/////////////////////
// Looping
/////////////////////

func (mp *MusicPlayer) SetLoopModeOff() {
	mp.loopMode = LoopOff
}

func (mp *MusicPlayer) SetLoopModeTrack() {
	mp.loopMode = LoopTrack
}

func (mp *MusicPlayer) SetLoopModeQueue() {
	mp.loopMode = LoopQueue
}

func (mp *MusicPlayer) IsLoopModeOff() bool {
	return mp.loopMode == LoopOff
}

func (mp *MusicPlayer) IsLoopModeTrack() bool {
	return mp.loopMode == LoopTrack
}

func (mp *MusicPlayer) IsLoopModeQueue() bool {
	return mp.loopMode == LoopQueue
}

/////////////////////
// Getters / Bools
/////////////////////

func (mp *MusicPlayer) GetLoadedTrack() (*lavalink.Track, error) {
	if !mp.HasLoadedTrack() {
		return nil, errors.New(util.NoLoadedTrack)
	}

	return mp.player.Track(), nil
}

func (mp *MusicPlayer) HasLoadedTrack() bool {
	return mp.player.Track() != nil
}

func (mp *MusicPlayer) GetPosition() lavalink.Duration {
	return mp.player.Position()
}

func (mp *MusicPlayer) IsPaused() bool {
	return mp.player.Paused()
}

/////////////////////
// Queue
/////////////////////

func (mp *MusicPlayer) GetQueue() []lavalink.Track {
	return mp.queue.PeekList()
}

func (mp *MusicPlayer) IsQueueEmpty() bool {
	return mp.queue.IsEmpty()
}

func (mp *MusicPlayer) GetQueueLength() int {
	return mp.queue.Length()
}

func (mp *MusicPlayer) GetQueueDuration() lavalink.Duration {
	return mp.queue.Duration()
}

func (mp *MusicPlayer) ClearQueue(num int) {
	if num <= 0 || num >= mp.queue.size {
		mp.queue.Clear()
		return
	}

	for i := 0; i < num; i++ {
		mp.queue.Dequeue()
	}
}

func (mp *MusicPlayer) ShuffleQueue() {
	mp.queue.Shuffle()
}

/////////////////////
// Helper Functions
/////////////////////

func (mp *MusicPlayer) newPlayer() {
	player := (*mp.lavalinkClient).Player(mp.guildID)

	player.Update(context.TODO(), lavalink.WithVolume(42))

	//player.AddListener(NewEventListener(&mp.Queue, &mp.LoopMode, &mp.GotDisconnected))
	mp.player = player
}

func (mp *MusicPlayer) recreatePlayer() {
	// Should be automatically destroyed in disgolink, upon leaving a voice channel.
	if mp.disconnected {
		mp.newPlayer()
		mp.disconnected = false
	}
}
