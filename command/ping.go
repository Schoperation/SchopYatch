package command

import (
	"schoperation/schopyatch/util"
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

func (cmd *PingCmd) Execute(deps CommandDependencies, opts ...string) error {
	err := util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Pong!")
	if err != nil {
		return err
	}

	return nil
}
