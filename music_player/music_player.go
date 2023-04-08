package music_player

import (
	"context"
	"errors"
	"schoperation/schopyatch/enum"
	"schoperation/schopyatch/util"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgolink/v2/disgolink"
	"github.com/disgoorg/disgolink/v2/lavalink"
	"github.com/disgoorg/snowflake/v2"
)

type MusicPlayer struct {
	guildID        snowflake.ID
	lavalinkClient *disgolink.Client
	player         disgolink.Player
	queue          MusicQueue
	searchResults  SearchResults
	loopMode       enum.LoopMode
	disconnected   bool
}

func NewMusicPlayer(guildId snowflake.ID, lavaLinkClient *disgolink.Client) *MusicPlayer {
	musicPlayer := MusicPlayer{
		guildID:        guildId,
		lavalinkClient: lavaLinkClient,
		queue:          NewMusicQueue(),
		searchResults:  NewSearchResults(),
		loopMode:       enum.LoopOff,
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
	if mp.hasLoadedTrack() {
		_, err := mp.Stop()
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

func (mp *MusicPlayer) ProcessQuery(query string) (enum.PlayerStatus, int, error) {
	var playerStatus enum.PlayerStatus
	var tracksQueued int
	var playerErr error

	(*mp.lavalinkClient).BestNode().LoadTracksHandler(context.TODO(), query, disgolink.NewResultHandler(
		func(track lavalink.Track) {
			playerStatus, playerErr = mp.Load(track)
			tracksQueued = 0
		},
		func(playlist lavalink.Playlist) {
			playerStatus, tracksQueued, playerErr = mp.LoadList(playlist.Tracks)
		},
		func(tracks []lavalink.Track) {
			mp.SetSearchResults(tracks)
			playerStatus = enum.StatusSuccess
			tracksQueued = 0
		},
		func() {
			playerStatus = enum.StatusFailed
			tracksQueued = 0
			playerErr = errors.New(util.NoResultsFound)
		},
		func(err error) {
			playerStatus = enum.StatusFailed
			tracksQueued = 0
			playerErr = err
		},
	))

	return playerStatus, tracksQueued, playerErr
}

/////////////////////
// Basic Player Cmds
/////////////////////

func (mp *MusicPlayer) Load(track lavalink.Track) (enum.PlayerStatus, error) {
	if mp.hasLoadedTrack() {
		mp.queue.Enqueue(track)
		return enum.StatusQueued, nil
	}

	err := mp.player.Update(context.TODO(), lavalink.WithTrack(track), lavalink.WithPaused(false))
	if err != nil {
		return enum.StatusFailed, err
	}

	return enum.StatusSuccess, nil
}

func (mp *MusicPlayer) LoadList(tracks []lavalink.Track) (enum.PlayerStatus, int, error) {
	if len(tracks) == 0 {
		return enum.StatusSuccess, 0, nil
	}

	if !mp.hasLoadedTrack() {
		err := mp.player.Update(context.TODO(), lavalink.WithTrack(tracks[0]), lavalink.WithPaused(false))
		if err != nil {
			return enum.StatusFailed, 0, err
		}

		if len(tracks) == 1 {
			return enum.StatusSuccess, 0, nil
		}

		tracks = tracks[1:]

		// Exception here just so we don't mention a one-hit-wonder when there were actually 2 tracks.
		if len(tracks) == 1 {
			mp.queue.Enqueue(tracks[0])
			return enum.StatusPlayingAndQueuedList, 1, nil
		}
	}

	mp.queue.EnqueueList(tracks)
	return enum.StatusQueuedList, len(tracks), nil
}

func (mp *MusicPlayer) Pause() (enum.PlayerStatus, error) {
	if !mp.hasLoadedTrack() {
		return enum.StatusFailed, errors.New(util.NoLoadedTrack)
	}

	if mp.player.Paused() {
		return enum.StatusAlreadyPaused, nil
	}

	err := mp.player.Update(context.TODO(), lavalink.WithPaused(true))
	if err != nil {
		return enum.StatusFailed, err
	}

	return enum.StatusSuccess, nil
}

func (mp *MusicPlayer) Unpause() (enum.PlayerStatus, error) {
	if !mp.hasLoadedTrack() {
		return enum.StatusFailed, errors.New(util.NoLoadedTrack)
	}

	if !mp.player.Paused() {
		return enum.StatusAlreadyUnpaused, nil
	}

	err := mp.player.Update(context.TODO(), lavalink.WithPaused(false))
	if err != nil {
		return enum.StatusFailed, err
	}

	return enum.StatusSuccess, nil
}

func (mp *MusicPlayer) Stop() (enum.PlayerStatus, error) {
	if !mp.hasLoadedTrack() {
		return enum.StatusFailed, errors.New(util.NoLoadedTrack)
	}

	err := mp.player.Update(context.TODO(), lavalink.WithNullTrack())
	if err != nil {
		return enum.StatusFailed, err
	}

	return enum.StatusSuccess, nil
}

func (mp *MusicPlayer) Seek(time lavalink.Duration) (enum.PlayerStatus, error) {
	if !mp.hasLoadedTrack() {
		return enum.StatusFailed, errors.New(util.NoLoadedTrack)
	}

	if time.Milliseconds() < 0 {
		return enum.StatusFailed, errors.New(util.NegativeDuration)
	}

	err := mp.player.Update(context.TODO(), lavalink.WithPosition(time))
	if err != nil {
		return enum.StatusFailed, err
	}

	return 0, nil
}

func (mp *MusicPlayer) Skip() (*lavalink.Track, error) {
	if !mp.hasLoadedTrack() {
		return nil, errors.New(util.NoLoadedTrack)
	}

	if mp.IsLoopModeQueue() {
		mp.queue.Enqueue(*mp.player.Track())
	}

	if mp.queue.IsEmpty() {
		err := mp.player.Update(context.TODO(), lavalink.WithNullTrack(), lavalink.WithPaused(false))
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	nextTrack := mp.queue.Dequeue()
	err := mp.player.Update(context.TODO(), lavalink.WithTrack(*nextTrack), lavalink.WithPaused(false))
	if err != nil {
		return nil, err
	}

	return nextTrack, nil
}

func (mp *MusicPlayer) SkipTo(index int) (*lavalink.Track, error) {
	if !mp.hasLoadedTrack() {
		return nil, errors.New(util.NoLoadedTrack)
	}

	if mp.queue.IsEmpty() {
		return nil, nil
	}

	if index < 0 || index >= mp.queue.Length() {
		return nil, errors.New(util.IndexOutOfBounds)
	}

	nextTrack := mp.queue.PeekAt(index)
	err := mp.player.Update(context.TODO(), lavalink.WithTrack(*nextTrack), lavalink.WithPaused(false))
	if err != nil {
		return nil, err
	}

	if mp.IsLoopModeQueue() {
		mp.queue.Enqueue(*mp.player.Track())
	}

	for i := 0; i < index+1; i++ {
		track := mp.queue.Dequeue()

		if mp.IsLoopModeQueue() {
			mp.queue.Enqueue(*track)
		}
	}

	return nextTrack, nil
}

/////////////////////
// Looping
/////////////////////

func (mp *MusicPlayer) SetLoopModeOff() {
	mp.loopMode = enum.LoopOff
}

func (mp *MusicPlayer) SetLoopModeTrack() {
	mp.loopMode = enum.LoopTrack
}

func (mp *MusicPlayer) SetLoopModeQueue() {
	mp.loopMode = enum.LoopQueue
}

func (mp *MusicPlayer) IsLoopModeOff() bool {
	return mp.loopMode == enum.LoopOff
}

func (mp *MusicPlayer) IsLoopModeTrack() bool {
	return mp.loopMode == enum.LoopTrack
}

func (mp *MusicPlayer) IsLoopModeQueue() bool {
	return mp.loopMode == enum.LoopQueue
}

/////////////////////
// Search Results
/////////////////////

func (mp *MusicPlayer) GetSearchResults() []lavalink.Track {
	return mp.searchResults.GetResults()
}

func (mp *MusicPlayer) GetSearchResult(index int) *lavalink.Track {
	return mp.searchResults.GetResult(index)
}

func (mp *MusicPlayer) SetSearchResults(tracks []lavalink.Track) {
	mp.searchResults.AddResults(tracks)
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
// Getters / Bools
/////////////////////

func (mp *MusicPlayer) GetLoadedTrack() (*lavalink.Track, error) {
	if !mp.hasLoadedTrack() {
		return nil, errors.New(util.NoLoadedTrack)
	}

	return mp.player.Track(), nil
}

func (mp *MusicPlayer) GetPosition() lavalink.Duration {
	return mp.player.Position()
}

/////////////////////
// Internal Helper Functions
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

func (mp *MusicPlayer) hasLoadedTrack() bool {
	return mp.player.Track() != nil
}
