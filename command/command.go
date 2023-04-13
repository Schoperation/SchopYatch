package command

import (
	"schoperation/schopyatch/enum"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgolink/v2/lavalink"
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/events"
)

type CommandDependencies struct {
	Client      *bot.Client
	Event       *events.MessageCreate
	MusicPlayer MusicPlayer
	Messenger   Messenger
	Prefix      string
}

type Command interface {
	GetName() string
	GetSummary() string
	GetDescription() string
	GetUsage() string
	GetAliases() []string
	IsVoiceOnlyCmd() bool
	Execute(deps CommandDependencies, opts ...string) error
}

type Messenger interface {
	SendSimpleMessage(msg string)
}

// TODO Begin implementing tests
type MusicPlayer interface {
	AddTrackToQueue(track lavalink.Track)
	ClearQueue(num int)
	GetLoadedTrack() (*lavalink.Track, error)
	GetPosition() lavalink.Duration
	GetQueue() []lavalink.Track
	GetQueueDuration() lavalink.Duration
	GetQueueLength() int
	GetSearchResult(index int) *lavalink.Track
	GetSearchResults() []lavalink.Track
	GetSearchResultsLength() int
	IsLoopModeOff() bool
	IsLoopModeQueue() bool
	IsLoopModeTrack() bool
	IsPaused() bool
	IsQueueEmpty() bool
	JoinVoiceChannel(botClient *bot.Client, userId snowflake.ID) error
	LeaveVoiceChannel(botClient *bot.Client, shouldReset bool) error
	Load(track lavalink.Track) (enum.PlayerStatus, error)
	LoadList(tracks []lavalink.Track) (enum.PlayerStatus, int, error)
	Pause() (enum.PlayerStatus, error)
	ProcessQuery(query string) (enum.PlayerStatus, *lavalink.Track, int, error)
	RemoveNextTrackFromQueue() (*lavalink.Track, error)
	RemoveTrackFromQueue(index int) (*lavalink.Track, error)
	Seek(time lavalink.Duration) (enum.PlayerStatus, error)
	SetLoopModeOff()
	SetLoopModeQueue()
	SetLoopModeTrack()
	SetSearchResults(tracks []lavalink.Track)
	ShuffleQueue()
	Skip() (*lavalink.Track, error)
	SkipTo(index int) (*lavalink.Track, error)
	Stop() (enum.PlayerStatus, error)
	Unpause() (enum.PlayerStatus, error)
}
