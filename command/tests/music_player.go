package command_tests

import (
	"schoperation/schopyatch/enum"
	"schoperation/schopyatch/music_player"
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgolink/v2/lavalink"
	"github.com/disgoorg/snowflake/v2"
)

type fakeMusicPlayer struct {
	guildID       snowflake.ID
	queue         music_player.MusicQueue
	searchResults music_player.SearchResults
	loopMode      enum.LoopMode
	disconnected  bool

	// Testing Only
	LoadedTrack     *lavalink.Track
	Paused          bool
	CurrentPosition lavalink.Duration
	StatusToReturn  enum.PlayerStatus
	ErrorsToReturn  map[string]error
	TracksQueued    int
}

func NewDefaultFakeMusicPlayer() fakeMusicPlayer {
	return fakeMusicPlayer{
		guildID:       snowflake.New(time.Now().UTC()),
		queue:         music_player.NewMusicQueue(),
		searchResults: music_player.NewSearchResults(),
		loopMode:      enum.LoopOff,
		disconnected:  false,
	}
}

func (fmp *fakeMusicPlayer) AddTrackToQueue(track lavalink.Track) {
	fmp.queue.Enqueue(track)
}

func (fmp *fakeMusicPlayer) ClearQueue(num int) {
	fmp.queue.Clear()
}

func (fmp *fakeMusicPlayer) GetLoadedTrack() (*lavalink.Track, error) {
	return fmp.LoadedTrack, fmp.ErrorsToReturn["GetLoadedTrack"]
}

func (fmp *fakeMusicPlayer) GetPosition() lavalink.Duration {
	return fmp.CurrentPosition
}

func (fmp *fakeMusicPlayer) GetQueue() []lavalink.Track {
	return fmp.queue.PeekList()
}

func (fmp *fakeMusicPlayer) GetQueueDuration() lavalink.Duration {
	return fmp.queue.Duration()
}

func (fmp *fakeMusicPlayer) GetQueueLength() int {
	return fmp.queue.Length()
}

func (fmp *fakeMusicPlayer) GetSearchResult(index int) *lavalink.Track {
	return fmp.searchResults.GetResult(index)
}

func (fmp *fakeMusicPlayer) GetSearchResults() []lavalink.Track {
	return fmp.searchResults.GetResults()
}

func (fmp *fakeMusicPlayer) GetSearchResultsLength() int {
	return fmp.searchResults.Length()
}

func (fmp *fakeMusicPlayer) IsLoopModeOff() bool {
	return fmp.loopMode == enum.LoopOff
}

func (fmp *fakeMusicPlayer) IsLoopModeQueue() bool {
	return fmp.loopMode == enum.LoopQueue
}

func (fmp *fakeMusicPlayer) IsLoopModeTrack() bool {
	return fmp.loopMode == enum.LoopTrack
}

func (fmp *fakeMusicPlayer) IsPaused() bool {
	return fmp.Paused
}

func (fmp *fakeMusicPlayer) IsQueueEmpty() bool {
	return fmp.queue.IsEmpty()
}

func (fmp *fakeMusicPlayer) JoinVoiceChannel(botClient *bot.Client, userId snowflake.ID) error {
	return fmp.ErrorsToReturn["JoinVoiceChannel"]
}

func (fmp *fakeMusicPlayer) LeaveVoiceChannel(botClient *bot.Client) error {
	return fmp.ErrorsToReturn["LeaveVoiceChannel"]
}

func (fmp *fakeMusicPlayer) Load(track lavalink.Track) (enum.PlayerStatus, error) {
	return fmp.StatusToReturn, fmp.ErrorsToReturn["Load"]
}

func (fmp *fakeMusicPlayer) LoadList(tracks []lavalink.Track) (enum.PlayerStatus, int, error) {
	return fmp.StatusToReturn, 0, fmp.ErrorsToReturn["LoadList"]
}

func (fmp *fakeMusicPlayer) Pause() (enum.PlayerStatus, error) {
	return fmp.StatusToReturn, fmp.ErrorsToReturn["Pause"]
}

func (fmp *fakeMusicPlayer) ProcessQuery(query string) (enum.PlayerStatus, *lavalink.Track, int, error) {
	return fmp.StatusToReturn, fmp.LoadedTrack, fmp.TracksQueued, fmp.ErrorsToReturn["ProcessQuery"]
}

func (fmp *fakeMusicPlayer) RemoveNextTrackFromQueue() (*lavalink.Track, error) {
	return fmp.queue.Dequeue(), fmp.ErrorsToReturn["RemoveNextTrackFromQueue"]
}

func (fmp *fakeMusicPlayer) RemoveTrackFromQueue(index int) (*lavalink.Track, error) {
	return fmp.queue.DequeueAt(index), fmp.ErrorsToReturn["RemoveTrackFromQueue"]
}

func (fmp *fakeMusicPlayer) Seek(time lavalink.Duration) (enum.PlayerStatus, error) {
	return fmp.StatusToReturn, fmp.ErrorsToReturn["Seek"]
}

func (fmp *fakeMusicPlayer) SetLoopModeOff() {
	fmp.loopMode = enum.LoopOff
}

func (fmp *fakeMusicPlayer) SetLoopModeQueue() {
	fmp.loopMode = enum.LoopQueue
}

func (fmp *fakeMusicPlayer) SetLoopModeTrack() {
	fmp.loopMode = enum.LoopTrack
}

func (fmp *fakeMusicPlayer) SetSearchResults(tracks []lavalink.Track) {
	fmp.searchResults.AddResults(tracks)
}

func (fmp *fakeMusicPlayer) ShuffleQueue() error {
	return fmp.ErrorsToReturn["ShuffleQueue"]
}

func (fmp *fakeMusicPlayer) Skip() (*lavalink.Track, error) {
	return fmp.LoadedTrack, fmp.ErrorsToReturn["Skip"]
}

func (fmp *fakeMusicPlayer) SkipTo(index int) (*lavalink.Track, error) {
	return fmp.LoadedTrack, fmp.ErrorsToReturn["SkipTo"]
}

func (fmp *fakeMusicPlayer) Stop() (enum.PlayerStatus, error) {
	return fmp.StatusToReturn, fmp.ErrorsToReturn["Stop"]
}

func (fmp *fakeMusicPlayer) Unpause() (enum.PlayerStatus, error) {
	return fmp.StatusToReturn, fmp.ErrorsToReturn["Unpause"]
}
