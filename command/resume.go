package command

import (
	"schoperation/schopyatch/enum"
	"schoperation/schopyatch/util"
)

type ResumeCmd struct {
	name        string
	summary     string
	description string
	usage       string
	aliases     []string
	voiceOnly   bool
}

func NewResumeCmd() Command {
	return &ResumeCmd{
		name:        "resume",
		summary:     "Resume the player",
		description: "This command simply resumes the player if it's paused.",
		usage:       "resume",
		aliases:     []string{"unpause"},
		voiceOnly:   true,
	}
}

func (cmd *ResumeCmd) GetName() string {
	return cmd.name
}

func (cmd *ResumeCmd) GetSummary() string {
	return cmd.summary
}

func (cmd *ResumeCmd) GetDescription() string {
	return cmd.description
}

func (cmd *ResumeCmd) GetUsage() string {
	return cmd.usage
}

func (cmd *ResumeCmd) GetAliases() []string {
	return cmd.aliases
}

func (cmd *ResumeCmd) IsVoiceOnlyCmd() bool {
	return cmd.voiceOnly
}

func (cmd *ResumeCmd) Execute(deps CommandDependencies, opts ...string) error {
	status, err := deps.MusicPlayer.Unpause()
	if err != nil && util.IsErrorMessage(err, util.NoLoadedTrack) {
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "No track's currently playing. Are we resuming a val sesh? Oh boy, it's 4 AM already, shame...")
		return nil
	} else if err != nil {
		return err
	}

	if status == enum.StatusAlreadyUnpaused {
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Already playing. Bruh where's your ears? Better yet, where's the thing between them?")
		return nil
	}

	util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Resuming.")
	return nil
}
