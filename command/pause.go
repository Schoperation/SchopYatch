package command

import (
	"fmt"
	"schoperation/schopyatch/util"
)

type PauseCmd struct {
	Name        string
	Summary     string
	Description string
	Usage       string
	Aliases     []string
}

func NewPauseCmd() Command {
	return &PauseCmd{
		Name:        "pause",
		Summary:     "Pause the player",
		Description: "This command simply pauses the player. Use resume, unpause, or play to resume the track.",
		Usage:       "pause",
		Aliases:     []string{"holdit"},
	}
}

func (cmd *PauseCmd) GetName() string {
	return cmd.Name
}

func (cmd *PauseCmd) GetSummary() string {
	return cmd.Summary
}

func (cmd *PauseCmd) GetDescription() string {
	return cmd.Description
}

func (cmd *PauseCmd) GetUsage() string {
	return cmd.Usage
}

func (cmd *PauseCmd) GetAliases() []string {
	return cmd.Aliases
}

func (cmd *PauseCmd) Execute(deps CommandDependencies, opts ...string) error {
	if deps.MusicPlayer.Player.PlayingTrack() == nil {
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "No track's currently playing. Would you like to pause time instead?")
		return nil
	}

	if deps.MusicPlayer.Player.Paused() {
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Already paused the music, man. Why such a party pooper?")
		return nil
	}

	err := deps.MusicPlayer.Player.Pause(true)
	if err != nil {
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Can't pause the player for some reason...")
		return err
	}

	util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, fmt.Sprintf("Paused. Use %sresume, %sunpause, or %splay to resume.", deps.Prefix, deps.Prefix, deps.Prefix))
	return nil
}
