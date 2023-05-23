package command

import (
	"fmt"
	"schoperation/schopyatch/enum"
	"schoperation/schopyatch/msg"
)

type PauseCmd struct {
	name        string
	group       string
	summary     string
	description string
	usage       string
	aliases     []string
	voiceOnly   bool
}

func NewPauseCmd() Command {
	return &PauseCmd{
		name:        "pause",
		group:       "player",
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

func (cmd *PauseCmd) GetGroup() string {
	return cmd.group
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
	if err != nil && msg.IsErrorMessage(err, msg.NoLoadedTrack) {
		deps.Messenger.SendSimpleMessage("No track's currently playing. Would you like to pause time instead?")
		return nil
	} else if err != nil {
		return err
	}

	if status == enum.StatusAlreadyPaused {
		deps.Messenger.SendSimpleMessage("Already paused the music, man. Why such a party pooper?")
		return nil
	}

	deps.Messenger.SendSimpleMessage(fmt.Sprintf("Paused. Use `%sresume`, `%sunpause`, or `%splay` to resume.", deps.Prefix, deps.Prefix, deps.Prefix))
	return nil
}
