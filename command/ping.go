package command

type PingCmd struct {
	name        string
	summary     string
	description string
	usage       string
	aliases     []string
	voiceOnly   bool
}

func NewPingCmd() Command {
	return &PingCmd{
		name:        "ping",
		summary:     "Pong!",
		description: "Plays ping pong. Brilliant, I know...",
		usage:       "ping",
		aliases:     []string{"pong"},
		voiceOnly:   false,
	}
}

func (cmd *PingCmd) GetName() string {
	return cmd.name
}

func (cmd *PingCmd) GetSummary() string {
	return cmd.summary
}

func (cmd *PingCmd) GetDescription() string {
	return cmd.description
}

func (cmd *PingCmd) GetUsage() string {
	return cmd.usage
}

func (cmd *PingCmd) GetAliases() []string {
	return cmd.aliases
}

func (cmd *PingCmd) IsVoiceOnlyCmd() bool {
	return cmd.voiceOnly
}

func (cmd *PingCmd) Execute(deps CommandDependencies, opts ...string) error {
	deps.Messenger.SendSimpleMessage("Pong!")
	return nil
}
