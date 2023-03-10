package command

import (
	"fmt"
	"schoperation/schopyatch/util"
	"strconv"
)

type RemoveCmd struct {
	name        string
	summary     string
	description string
	usage       string
	aliases     []string
	voiceOnly   bool
}

func NewRemoveCmd() Command {
	return &RemoveCmd{
		name:        "remove",
		summary:     "Remove a track from the queue",
		description: "This command simply removes a specified track number from the queue.",
		usage:       "remove <number>",
		aliases:     []string{"r", "delete"},
		voiceOnly:   true,
	}
}

func (cmd *RemoveCmd) GetName() string {
	return cmd.name
}

func (cmd *RemoveCmd) GetSummary() string {
	return cmd.summary
}

func (cmd *RemoveCmd) GetDescription() string {
	return cmd.description
}

func (cmd *RemoveCmd) GetUsage() string {
	return cmd.usage
}

func (cmd *RemoveCmd) GetAliases() []string {
	return cmd.aliases
}
func (cmd *RemoveCmd) IsVoiceOnlyCmd() bool {
	return cmd.voiceOnly
}

func (cmd *RemoveCmd) Execute(deps CommandDependencies, opts ...string) error {
	if len(opts) == 0 {
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, fmt.Sprintf("Need to specify a number to remove. E.g. `%sremove 4` to remove the 4th track. Try `%squeue` to see the numbers.", deps.Prefix, deps.Prefix))
		return nil
	}

	num, err := strconv.Atoi(opts[0])
	if err != nil {
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Dude that's some voodoo... we need a number.")
		return err
	}

	if deps.MusicPlayer.Queue.IsEmpty() {
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Nothing to remove. The glass, this time, is fully empty.")
		return nil
	}

	if num > deps.MusicPlayer.Queue.Length() {
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, fmt.Sprintf("That index doesn't exist. The current length is %d.", deps.MusicPlayer.Queue.Length()))
		return nil
	}

	track := deps.MusicPlayer.Queue.DequeueAtIndex(num - 1)
	util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, fmt.Sprintf("Removed *%s* by **%s** from the queue.", (*track).Info().Title, (*track).Info().Author))
	return nil
}
