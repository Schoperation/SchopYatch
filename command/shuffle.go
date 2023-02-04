package command

import "schoperation/schopyatch/util"

type ShuffleCmd struct {
	name        string
	summary     string
	description string
	usage       string
	aliases     []string
}

func NewShuffleCmd() Command {
	return &ShuffleCmd{
		name:        "shuffle",
		summary:     "Shuffle the queue",
		description: "This command simply shuffles the current queue.",
		usage:       "shuffle",
		aliases:     []string{"riffle"},
	}
}

func (cmd *ShuffleCmd) GetName() string {
	return cmd.name
}

func (cmd *ShuffleCmd) GetSummary() string {
	return cmd.summary
}

func (cmd *ShuffleCmd) GetDescription() string {
	return cmd.description
}

func (cmd *ShuffleCmd) GetUsage() string {
	return cmd.usage
}

func (cmd *ShuffleCmd) GetAliases() []string {
	return cmd.aliases
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
