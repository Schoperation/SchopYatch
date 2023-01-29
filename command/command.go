package command

import (
	"schoperation/schopyatch/util"

	"github.com/disgoorg/disgo/bot"

	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgolink/disgolink"
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
	Client   *bot.Client
	Lavalink *disgolink.Link
	Event    *events.MessageCreate
	Queue    *util.MusicQueue
	Prefix   string
	LoopMode util.LoopMode
}
