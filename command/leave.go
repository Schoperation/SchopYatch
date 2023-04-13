package command

type LeaveCmd struct {
	name        string
	summary     string
	description string
	usage       string
	aliases     []string
	voiceOnly   bool
}

func NewLeaveCmd() Command {
	return &LeaveCmd{
		name:        "leave",
		summary:     "Make the bot leave a voice channel",
		description: "Upon running, the bot will leave the user's voice channel. Kindly. Add \"reset\" to clear the queue, search results, and reset looping as well.",
		usage:       "leave [reset]",
		aliases:     []string{"fuckoff", "disconnect"},
		voiceOnly:   true,
	}
}

func (cmd *LeaveCmd) GetName() string {
	return cmd.name
}

func (cmd *LeaveCmd) GetSummary() string {
	return cmd.summary
}

func (cmd *LeaveCmd) GetDescription() string {
	return cmd.description
}

func (cmd *LeaveCmd) GetUsage() string {
	return cmd.usage
}

func (cmd *LeaveCmd) GetAliases() []string {
	return cmd.aliases
}

func (cmd *LeaveCmd) IsVoiceOnlyCmd() bool {
	return cmd.voiceOnly
}

func (cmd *LeaveCmd) Execute(deps CommandDependencies, opts ...string) error {
	if len(opts) > 0 {
		if opts[0] == "reset" {
			err := deps.MusicPlayer.LeaveVoiceChannel(deps.Client, true)
			if err != nil {
				return err
			}

			return nil
		}

		deps.Messenger.SendSimpleMessage("You want me to leave what? Try either leaving it blank or `reset`.")
		return nil
	}

	err := deps.MusicPlayer.LeaveVoiceChannel(deps.Client, false)
	return err
}
