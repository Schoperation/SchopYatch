package command

import (
	"fmt"
	"strconv"
)

type ClearCmd struct {
	name        string
	summary     string
	description string
	usage       string
	aliases     []string
	voiceOnly   bool
}

func NewClearCmd() Command {
	return &ClearCmd{
		name:        "clear",
		summary:     "Clear the queue",
		description: "This command simply clears the queue. Has an additional, optional parameter to clear only the first num or so entries. Will not affect the track currently playing; use skip for that.",
		usage:       "clear [num]",
		aliases:     []string{"c", "empty", "clearqueue", "clearlist"},
		voiceOnly:   true,
	}
}

func (cmd *ClearCmd) GetName() string {
	return cmd.name
}

func (cmd *ClearCmd) GetSummary() string {
	return cmd.summary
}

func (cmd *ClearCmd) GetDescription() string {
	return cmd.description
}

func (cmd *ClearCmd) GetUsage() string {
	return cmd.usage
}

func (cmd *ClearCmd) GetAliases() []string {
	return cmd.aliases
}

func (cmd *ClearCmd) IsVoiceOnlyCmd() bool {
	return cmd.voiceOnly
}

func (cmd *ClearCmd) Execute(deps CommandDependencies, opts ...string) error {
	if deps.MusicPlayer.IsQueueEmpty() {
		deps.Messenger.SendSimpleMessage("Already clear. Were you hoping for a funny error? Same")
		return nil
	}

	if len(opts) == 0 {
		deps.MusicPlayer.ClearQueue(0)
		deps.Messenger.SendSimpleMessage("Cleared the queue.")
		return nil
	}

	num, err := strconv.Atoi(opts[0])
	if err != nil {
		deps.Messenger.SendSimpleMessage("Dude, that was not a number...")
		return nil
	}

	deps.MusicPlayer.ClearQueue(num)
	deps.Messenger.SendSimpleMessage(fmt.Sprintf("Cleared the first %d tracks from the queue.", num))
	return nil
}
