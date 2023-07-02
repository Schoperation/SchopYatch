package command

import (
	"schoperation/schopyatch/enum"
	"schoperation/schopyatch/music_player"
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/events"
)

type CommandDependencies struct {
	Event       *events.MessageCreate
	MusicPlayer MusicPlayer
	Messenger   Messenger
	Prefix      string
	Version     string
}

type Command interface {
	GetName() string
	GetGroup() string
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

// TODO Break this up into smaller interfaces, per command
type MusicPlayer interface {
	AddTrackToQueue(track music_player.Track)
	ClearQueue(num int)
	GetLoadedTrack() (*music_player.Track, error)
	GetPosition() time.Duration
	GetQueue() []music_player.Track
	GetQueueDuration() time.Duration
	GetQueueLength() int
	GetSearchResult(index int) *music_player.Track
	GetSearchResults() []music_player.Track
	GetSearchResultsLength() int
	IsLoopModeOff() bool
	IsLoopModeQueue() bool
	IsLoopModeTrack() bool
	IsPaused() bool
	IsQueueEmpty() bool
	JoinVoiceChannel(botClient *bot.Client, userId snowflake.ID) error
	LeaveVoiceChannel(botClient *bot.Client) error
	Load(track music_player.Track) (enum.PlayerStatus, error)
	LoadList(tracks []music_player.Track) (enum.PlayerStatus, int, error)
	Pause() (enum.PlayerStatus, error)
	ProcessQuery(query string) (enum.PlayerStatus, *music_player.Track, int, error)
	RemoveNextTrackFromQueue() (*music_player.Track, error)
	RemoveTrackFromQueue(index int) (*music_player.Track, error)
	Seek(time time.Duration) (enum.PlayerStatus, error)
	SetLoopModeOff()
	SetLoopModeQueue()
	SetLoopModeTrack()
	SetSearchResults(tracks []music_player.Track)
	ShuffleQueue() error
	Skip() (*music_player.Track, error)
	SkipTo(index int) (*music_player.Track, error)
	Stop() (enum.PlayerStatus, error)
	Unpause() (enum.PlayerStatus, error)
}
