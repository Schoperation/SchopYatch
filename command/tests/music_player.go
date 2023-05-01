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
	LoadedTrack    *lavalink.Track
	StatusToReturn enum.PlayerStatus
	ErrorToReturn  error
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
	return fmp.LoadedTrack, nil
}

func (fmp *fakeMusicPlayer) GetPosition() lavalink.Duration {
	return lavalink.Day
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
	return false
}

func (fmp *fakeMusicPlayer) IsQueueEmpty() bool {
	return fmp.queue.IsEmpty()
}

func (fmp *fakeMusicPlayer) JoinVoiceChannel(botClient *bot.Client, userId snowflake.ID) error {
	return nil
}

func (fmp *fakeMusicPlayer) LeaveVoiceChannel(botClient *bot.Client) error {
	return nil
}

func (fmp *fakeMusicPlayer) Load(track lavalink.Track) (enum.PlayerStatus, error) {
	return enum.StatusSuccess, nil
}

func (fmp *fakeMusicPlayer) LoadList(tracks []lavalink.Track) (enum.PlayerStatus, int, error) {
	return enum.StatusSuccess, 0, nil
}

func (fmp *fakeMusicPlayer) Pause() (enum.PlayerStatus, error) {
	return fmp.StatusToReturn, fmp.ErrorToReturn
}

func (fmp *fakeMusicPlayer) ProcessQuery(query string) (enum.PlayerStatus, *lavalink.Track, int, error) {
	return enum.StatusSuccess, nil, 0, nil
}

func (fmp *fakeMusicPlayer) RemoveNextTrackFromQueue() (*lavalink.Track, error) {
	return fmp.queue.Dequeue(), nil
}

func (fmp *fakeMusicPlayer) RemoveTrackFromQueue(index int) (*lavalink.Track, error) {
	return fmp.queue.DequeueAt(index), nil
}

func (fmp *fakeMusicPlayer) Seek(time lavalink.Duration) (enum.PlayerStatus, error) {
	return enum.StatusSuccess, nil
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
	return fmp.ErrorToReturn
}

func (fmp *fakeMusicPlayer) Skip() (*lavalink.Track, error) {
	return fmp.LoadedTrack, fmp.ErrorToReturn
}

func (fmp *fakeMusicPlayer) SkipTo(index int) (*lavalink.Track, error) {
	return fmp.LoadedTrack, fmp.ErrorToReturn
}

func (fmp *fakeMusicPlayer) Stop() (enum.PlayerStatus, error) {
	return fmp.StatusToReturn, fmp.ErrorToReturn
}

func (fmp *fakeMusicPlayer) Unpause() (enum.PlayerStatus, error) {
	return fmp.StatusToReturn, fmp.ErrorToReturn
}
