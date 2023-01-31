package command

import (
	"fmt"
	"schoperation/schopyatch/util"
	"strconv"
)

type ClearCmd struct {
	Name        string
	Summary     string
	Description string
	Usage       string
	Aliases     []string
}

func NewClearCmd() Command {
	return &ClearCmd{
		Name:        "clear",
		Summary:     "Clear the queue",
		Description: "This command simply clears the queue. Has an additional, optional parameter to clear only the first num or so entries. Will not affect the track currently playing; use skip for that.",
		Usage:       "clear [num]",
		Aliases:     []string{"c", "empty", "clearqueue", "clearlist"},
	}
}

func (cmd *ClearCmd) GetName() string {
	return cmd.Name
}

func (cmd *ClearCmd) GetSummary() string {
	return cmd.Summary
}

func (cmd *ClearCmd) GetDescription() string {
	return cmd.Description
}

func (cmd *ClearCmd) GetUsage() string {
	return cmd.Usage
}

func (cmd *ClearCmd) GetAliases() []string {
	return cmd.Aliases
}

func (cmd *ClearCmd) Execute(deps CommandDependencies, opts ...string) error {
	if deps.MusicPlayer.Queue.IsEmpty() {
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Already clear. Were you hoping for a funny error? Same")
		return nil
	}

	if len(opts) == 0 {
		deps.MusicPlayer.Queue.Clear()
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Cleared the queue.")
		return nil
	}

	num, err := strconv.Atoi(opts[0])
	if err != nil {
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Dude that was not a number...")
		return err
	}

	for i := 0; i < num; i++ {
		_ = deps.MusicPlayer.Queue.Dequeue()
	}

	util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, fmt.Sprintf("Cleared the first %d tracks from the queue.", num))
	return nil
}
