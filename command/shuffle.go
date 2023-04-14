package command

import "schoperation/schopyatch/msg"

type ShuffleCmd struct {
	name        string
	summary     string
	description string
	usage       string
	aliases     []string
	voiceOnly   bool
}

func NewShuffleCmd() Command {
	return &ShuffleCmd{
		name:        "shuffle",
		summary:     "Shuffle the queue",
		description: "This command simply shuffles the current queue.",
		usage:       "shuffle",
		aliases:     []string{"riffle"},
		voiceOnly:   true,
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

func (cmd *ShuffleCmd) IsVoiceOnlyCmd() bool {
	return cmd.voiceOnly
}

func (cmd *ShuffleCmd) Execute(deps CommandDependencies, opts ...string) error {
	err := deps.MusicPlayer.ShuffleQueue()
	if err != nil {
		if msg.IsErrorMessage(err, msg.QueueIsEmpty) {
			deps.Messenger.SendSimpleMessage("Nothing to shuffle. How else am I gonna show off my riffles?")
			return nil
		}

		return err
	}

	deps.Messenger.SendSimpleMessage("Shuffled the queue.")
	return nil
}
