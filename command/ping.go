package command

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

type PingCmd struct {
	Name        string
	Description string
}

func NewPingCmd() Command {
	return &PingCmd{
		Name:        "ping",
		Description: "Plays ping pong. Brilliant, I know...",
	}
}

func (cmd *PingCmd) GetName() string {
	return cmd.Name
}

func (cmd *PingCmd) GetDescription() string {
	return cmd.Description
}

func (cmd *PingCmd) Execute(e *events.MessageCreate, opts ...string) error {
	_, err := e.Client().Rest().CreateMessage(e.ChannelID, discord.NewMessageCreateBuilder().SetContent("Pong!").Build())
	if err != nil {
		return err
	}

	return nil
}
