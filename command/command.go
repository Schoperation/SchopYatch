package command

import (
	"github.com/disgoorg/disgo/bot"

	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgolink/disgolink"
)

type Command interface {
	GetName() string
	GetDescription() string
	Execute(deps CommandDependencies, opts ...string) error
}

type CommandDependencies struct {
	Client   *bot.Client
	Lavalink *disgolink.Link
	Event    *events.MessageCreate
}
