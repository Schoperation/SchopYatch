package command

import (
	"schoperation/schopyatch/util"
)

type PingCmd struct {
	Name        string
	Summary     string
	Description string
	Usage       string
	Aliases     []string
}

func NewPingCmd() Command {
	return &PingCmd{
		Name:        "ping",
		Summary:     "Pong!",
		Description: "Plays ping pong. Brilliant, I know...",
		Usage:       "ping",
		Aliases:     []string{"pong"},
	}
}

func (cmd *PingCmd) GetName() string {
	return cmd.Name
}

func (cmd *PingCmd) GetSummary() string {
	return cmd.Summary
}

func (cmd *PingCmd) GetDescription() string {
	return cmd.Description
}

func (cmd *PingCmd) GetUsage() string {
	return cmd.Usage
}

func (cmd *PingCmd) GetAliases() []string {
	return cmd.Aliases
}

func (cmd *PingCmd) Execute(deps CommandDependencies, opts ...string) error {
	util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Pong!")
	return nil
}
