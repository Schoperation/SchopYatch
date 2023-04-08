package command

import (
	"fmt"
	"schoperation/schopyatch/enum"
	"schoperation/schopyatch/util"
)

type PauseCmd struct {
	name        string
	summary     string
	description string
	usage       string
	aliases     []string
	voiceOnly   bool
}

func NewPauseCmd() Command {
	return &PauseCmd{
		name:        "pause",
		summary:     "Pause the player",
		description: "This command simply pauses the player. Use resume, unpause, or play to resume the track.",
		usage:       "pause",
		aliases:     []string{"holdit"},
		voiceOnly:   true,
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

func (cmd *PauseCmd) IsVoiceOnlyCmd() bool {
	return cmd.voiceOnly
}

func (cmd *PauseCmd) Execute(deps CommandDependencies, opts ...string) error {
	status, err := deps.MusicPlayer.Pause()
	if err != nil && util.IsErrorMessage(err, util.NoLoadedTrack) {
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "No track's currently playing. Would you like to pause time instead?")
		return nil
	} else if err != nil {
		return err
	}

	if status == enum.StatusAlreadyPaused {
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Already paused the music, man. Why such a party pooper?")
		return nil
	}

	util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, fmt.Sprintf("Paused. Use `%sresume`, `%sunpause`, or `%splay` to resume.", deps.Prefix, deps.Prefix, deps.Prefix))
	return nil
}
