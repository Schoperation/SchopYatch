package command

import (
	"schoperation/schopyatch/musicplayer"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgolink/disgolink"

	"github.com/disgoorg/disgo/events"
)

type Command interface {
	GetName() string
	GetSummary() string
	GetDescription() string
	GetUsage() string
	GetAliases() []string
	Execute(deps CommandDependencies, opts ...string) error
}

type CommandDependencies struct {
	Client      *bot.Client
	Event       *events.MessageCreate
	MusicPlayer *musicplayer.MusicPlayer
	Lavalink    *disgolink.Link
	Prefix      string
}
