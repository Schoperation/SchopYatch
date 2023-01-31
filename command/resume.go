package command

import "schoperation/schopyatch/util"

type ResumeCmd struct {
	Name        string
	Summary     string
	Description string
	Usage       string
	Aliases     []string
}

func NewResumeCmd() Command {
	return &ResumeCmd{
		Name:        "resume",
		Summary:     "Resume the player",
		Description: "This command simply resumes the player if it's paused.",
		Usage:       "resume",
		Aliases:     []string{"unpause"},
	}
}

func (cmd *ResumeCmd) GetName() string {
	return cmd.Name
}

func (cmd *ResumeCmd) GetSummary() string {
	return cmd.Summary
}

func (cmd *ResumeCmd) GetDescription() string {
	return cmd.Description
}

func (cmd *ResumeCmd) GetUsage() string {
	return cmd.Usage
}

func (cmd *ResumeCmd) GetAliases() []string {
	return cmd.Aliases
}

func (cmd *ResumeCmd) Execute(deps CommandDependencies, opts ...string) error {
	err := resume(deps)
	return err
}

func resume(deps CommandDependencies) error {
	if deps.MusicPlayer.Player.PlayingTrack() == nil {
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "No track's currently playing. Are we resuming a val sesh? Oh boy, it's 4 AM already, shame...")
		return nil
	}

	if !deps.MusicPlayer.Player.Paused() {
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Already playing. Bruh where's your ears? Better yet, where's the thing between them?")
		return nil
	}

	err := deps.MusicPlayer.Player.Pause(false)
	if err != nil {
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Can't resume the player for some reason...")
		return err
	}

	util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Resuming.")
	return nil
}
