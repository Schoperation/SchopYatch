package command

import (
	"schoperation/schopyatch/music_player"

	"github.com/disgoorg/disgo/bot"

	"github.com/disgoorg/disgo/events"
)

type Command interface {
	GetName() string
	GetSummary() string
	GetDescription() string
	GetUsage() string
	GetAliases() []string
	IsVoiceOnlyCmd() bool
	Execute(deps CommandDependencies, opts ...string) error
}

type CommandDependencies struct {
	Client      *bot.Client
	Event       *events.MessageCreate
	MusicPlayer *music_player.MusicPlayer
	Prefix      string
}
