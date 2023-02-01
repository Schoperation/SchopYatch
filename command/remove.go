package command

import (
	"fmt"
	"schoperation/schopyatch/util"
	"strconv"
)

type RemoveCmd struct {
	Name        string
	Summary     string
	Description string
	Usage       string
	Aliases     []string
}

func NewRemoveCmd() Command {
	return &RemoveCmd{
		Name:        "remove",
		Summary:     "Remove a track from the queue",
		Description: "This command simply removes a specified track number from the queue.",
		Usage:       "remove <number>",
		Aliases:     []string{"r", "delete"},
	}
}

func (cmd *RemoveCmd) GetName() string {
	return cmd.Name
}

func (cmd *RemoveCmd) GetSummary() string {
	return cmd.Summary
}

func (cmd *RemoveCmd) GetDescription() string {
	return cmd.Description
}

func (cmd *RemoveCmd) GetUsage() string {
	return cmd.Usage
}

func (cmd *RemoveCmd) GetAliases() []string {
	return cmd.Aliases
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
