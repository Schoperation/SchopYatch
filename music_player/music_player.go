package music_player

import (
	"context"
	"errors"
	"schoperation/schopyatch/enum"
	"schoperation/schopyatch/msg"
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgolink/v2/disgolink"
	"github.com/disgoorg/disgolink/v2/lavalink"
	"github.com/disgoorg/snowflake/v2"
)

// TODO add client interface to allow testing?

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
		return errors.New(msg.VoiceStateNotFound)
	}

	err := (*botClient).UpdateVoiceState(context.TODO(), mp.guildID, voiceState.ChannelID, false, true)
	if err != nil {
		return err
	}

	mp.recreatePlayer()
	return nil
}

func (mp *MusicPlayer) LeaveVoiceChannel(botClient *bot.Client) error {
	if mp.hasLoadedTrack() {
		err := mp.player.Update(context.TODO(), lavalink.WithNullTrack())
		if err != nil {
			return err
		}
	}

	err := (*botClient).UpdateVoiceState(context.TODO(), mp.guildID, nil, false, false)
	if err != nil {
		return err
	}

	mp.disconnected = true
	return nil
}

func (mp *MusicPlayer) ProcessQuery(query string) (enum.PlayerStatus, *Track, int, error) {
	var newTrack Track
	var playerStatus enum.PlayerStatus
	var tracksQueued int
	var playerErr error

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	(*mp.lavalinkClient).BestNode().LoadTracksHandler(ctx, query, disgolink.NewResultHandler(
		func(track lavalink.Track) {
			newTrack = toTrack(track)
			playerStatus, playerErr = mp.Load(toTrack(track))
			tracksQueued = 0
		},
		func(playlist lavalink.Playlist) {
			if len(playlist.Tracks) > 0 {
				newTrack = toTrack(playlist.Tracks[0])
			}

			playerStatus, tracksQueued, playerErr = mp.LoadList(toTracks(playlist.Tracks))
		},
		func(tracks []lavalink.Track) {
			mp.SetSearchResults(toTracks(tracks))
			playerStatus = enum.StatusSearchSuccess
			tracksQueued = 0
		},
		func() {
			playerStatus = enum.StatusFailed
			tracksQueued = 0
			playerErr = errors.New(msg.NoResultsFound)
		},
		func(err error) {
			playerStatus = enum.StatusFailed
			tracksQueued = 0
			playerErr = err
		},
	))

	return playerStatus, &newTrack, tracksQueued, playerErr
}

/////////////////////
// Basic Player Cmds
/////////////////////

func (mp *MusicPlayer) Load(track Track) (enum.PlayerStatus, error) {
	if mp.hasLoadedTrack() {
		mp.queue.Enqueue(track)
		return enum.StatusQueued, nil
	}

	err := mp.player.Update(context.TODO(), lavalink.WithTrack(track.ToLavalinkTrack()), lavalink.WithPaused(false))
	if err != nil {
		return enum.StatusFailed, err
	}

	return enum.StatusSuccess, nil
}

func (mp *MusicPlayer) LoadList(tracks []Track) (enum.PlayerStatus, int, error) {
	if len(tracks) == 0 {
		return enum.StatusQueuedList, 0, nil
	}

	if !mp.hasLoadedTrack() {
		err := mp.player.Update(context.TODO(), lavalink.WithTrack(tracks[0].ToLavalinkTrack()), lavalink.WithPaused(false))
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

// TODO make a play function that overrides currently playing song

func (mp *MusicPlayer) Pause() (enum.PlayerStatus, error) {
	if !mp.hasLoadedTrack() {
		return enum.StatusFailed, errors.New(msg.NoLoadedTrack)
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
		return enum.StatusFailed, errors.New(msg.NoLoadedTrack)
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
		return enum.StatusFailed, errors.New(msg.NoLoadedTrack)
	}

	err := mp.player.Update(context.TODO(), lavalink.WithNullTrack())
	if err != nil {
		return enum.StatusFailed, err
	}

	return enum.StatusSuccess, nil
}

func (mp *MusicPlayer) Seek(time time.Duration) (enum.PlayerStatus, error) {
	if !mp.hasLoadedTrack() {
		return enum.StatusFailed, errors.New(msg.NoLoadedTrack)
	}

	if time.Milliseconds() < 0 {
		return enum.StatusFailed, errors.New(msg.NegativeDuration)
	}

	if time >= toTimeDuration(mp.player.Track().Info.Length) {
		return enum.StatusFailed, errors.New(msg.IndexOutOfBounds)
	}

	err := mp.player.Update(context.TODO(), lavalink.WithPosition(toLavalinkDuration(time)))
	if err != nil {
		return enum.StatusFailed, err
	}

	return 0, nil
}

func (mp *MusicPlayer) Skip() (*Track, error) {
	if !mp.hasLoadedTrack() {
		return nil, errors.New(msg.NoLoadedTrack)
	}

	if mp.IsLoopModeQueue() {
		mp.queue.Enqueue(toTrack(*mp.player.Track()))
	}

	if mp.queue.IsEmpty() {
		err := mp.player.Update(context.TODO(), lavalink.WithNullTrack(), lavalink.WithPaused(false))
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	nextTrack := mp.queue.Dequeue()
	err := mp.player.Update(context.TODO(), lavalink.WithTrack(nextTrack.ToLavalinkTrack()), lavalink.WithPaused(false))
	if err != nil {
		return nil, err
	}

	return nextTrack, nil
}

func (mp *MusicPlayer) SkipTo(index int) (*Track, error) {
	if !mp.hasLoadedTrack() {
		return nil, errors.New(msg.NoLoadedTrack)
	}

	if mp.queue.IsEmpty() {
		err := mp.player.Update(context.TODO(), lavalink.WithNullTrack(), lavalink.WithPaused(false))
		if err != nil {
			return nil, err
		}

		return nil, nil
	}

	if index < 0 || index >= mp.queue.Length() {
		return nil, errors.New(msg.IndexOutOfBounds)
	}

	nextTrack := mp.queue.PeekAt(index)
	err := mp.player.Update(context.TODO(), lavalink.WithTrack(nextTrack.ToLavalinkTrack()), lavalink.WithPaused(false))
	if err != nil {
		return nil, err
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

func (mp *MusicPlayer) GetSearchResults() []Track {
	return mp.searchResults.GetResults()
}

func (mp *MusicPlayer) GetSearchResult(index int) *Track {
	return mp.searchResults.GetResult(index)
}

func (mp *MusicPlayer) GetSearchResultsLength() int {
	return mp.searchResults.length
}

func (mp *MusicPlayer) SetSearchResults(tracks []Track) {
	mp.searchResults.Clear()
	mp.searchResults.AddResults(tracks)
}

/////////////////////
// Queue
/////////////////////

func (mp *MusicPlayer) GetQueue() []Track {
	return mp.queue.PeekList()
}

func (mp *MusicPlayer) IsQueueEmpty() bool {
	return mp.queue.IsEmpty()
}

func (mp *MusicPlayer) GetQueueLength() int {
	return mp.queue.Length()
}

func (mp *MusicPlayer) GetQueueDuration() time.Duration {
	return mp.queue.Duration()
}

func (mp *MusicPlayer) ClearQueue(num int) {
	if num <= 0 || num >= mp.queue.Length() {
		mp.queue.Clear()
		return
	}

	for i := 0; i < num; i++ {
		mp.queue.Dequeue()
	}
}

func (mp *MusicPlayer) AddTrackToQueue(track Track) {
	mp.queue.Enqueue(track)
}

func (mp *MusicPlayer) RemoveNextTrackFromQueue() (*Track, error) {
	return mp.RemoveTrackFromQueue(0)
}

func (mp *MusicPlayer) RemoveTrackFromQueue(index int) (*Track, error) {
	if mp.queue.IsEmpty() {
		return nil, errors.New(msg.QueueIsEmpty)
	}

	if index < 0 || index >= mp.queue.Length() {
		return nil, errors.New(msg.IndexOutOfBounds)
	}

	track := mp.queue.DequeueAt(index)
	return track, nil
}

func (mp *MusicPlayer) ShuffleQueue() error {
	if mp.queue.IsEmpty() {
		return errors.New(msg.QueueIsEmpty)
	}

	mp.queue.Shuffle()
	return nil
}

/////////////////////
// Getters / Bools
/////////////////////

func (mp *MusicPlayer) GetLoadedTrack() (*Track, error) {
	if !mp.hasLoadedTrack() {
		return nil, errors.New(msg.NoLoadedTrack)
	}

	track := toTrack(*mp.player.Track())
	return &track, nil
}

func (mp *MusicPlayer) GetPosition() time.Duration {
	return toTimeDuration(mp.player.Position())
}

func (mp *MusicPlayer) IsPaused() bool {
	return mp.player.Paused()
}

/////////////////////
// Internal Helper Functions
/////////////////////

func (mp *MusicPlayer) newPlayer() {
	player := (*mp.lavalinkClient).Player(mp.guildID)

	player.Update(context.TODO(), lavalink.WithVolume(42))

	mp.player = player
}

func (mp *MusicPlayer) recreatePlayer() {
	// Player should be automatically destroyed in disgolink upon leaving a voice channel.
	if mp.disconnected {
		mp.newPlayer()
		mp.disconnected = false
	}
}

func (mp *MusicPlayer) hasLoadedTrack() bool {
	return mp.player.Track() != nil
}
