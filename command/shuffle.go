package command

import "schoperation/schopyatch/util"

type ShuffleCmd struct {
	Name        string
	Summary     string
	Description string
	Usage       string
	Aliases     []string
}

func NewShuffleCmd() Command {
	return &ShuffleCmd{
		Name:        "shuffle",
		Summary:     "Shuffle the queue",
		Description: "This command simply shuffles the current queue.",
		Usage:       "shuffle",
		Aliases:     []string{"riffle"},
	}
}

func (cmd *ShuffleCmd) GetName() string {
	return cmd.Name
}

func (cmd *ShuffleCmd) GetSummary() string {
	return cmd.Summary
}

func (cmd *ShuffleCmd) GetDescription() string {
	return cmd.Description
}

func (cmd *ShuffleCmd) GetUsage() string {
	return cmd.Usage
}

func (cmd *ShuffleCmd) GetAliases() []string {
	return cmd.Aliases
}

func (cmd *ShuffleCmd) Execute(deps CommandDependencies, opts ...string) error {
	if deps.MusicPlayer.Queue.IsEmpty() {
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Nothing to shuffle. How else am I gonna show off my riffles?")
		return nil
	}

	deps.MusicPlayer.Queue.Shuffle()
	util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Shuffled the queue.")
	return nil
}
