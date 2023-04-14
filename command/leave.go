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
		description: "Upon running, the bot will leave the user's voice channel, clearing the queue, search results, and everything else.",
		usage:       "leave",
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
	err := deps.MusicPlayer.LeaveVoiceChannel(deps.Client)
	return err
}
