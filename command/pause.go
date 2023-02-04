package command

import (
	"fmt"
	"schoperation/schopyatch/util"
)

type PauseCmd struct {
	name        string
	summary     string
	description string
	usage       string
	aliases     []string
}

func NewPauseCmd() Command {
	return &PauseCmd{
		name:        "pause",
		summary:     "Pause the player",
		description: "This command simply pauses the player. Use resume, unpause, or play to resume the track.",
		usage:       "pause",
		aliases:     []string{"holdit"},
	}
}

func (cmd *PauseCmd) GetName() string {
	return cmd.name
}

func (cmd *PauseCmd) GetSummary() string {
	return cmd.summary
}

func (cmd *PauseCmd) GetDescription() string {
	return cmd.description
}

func (cmd *PauseCmd) GetUsage() string {
	return cmd.usage
}

func (cmd *PauseCmd) GetAliases() []string {
	return cmd.aliases
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
