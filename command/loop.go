package command

type LoopCmd struct {
	name        string
	group       string
	summary     string
	description string
	usage       string
	aliases     []string
	voiceOnly   bool
}

func NewLoopCmd() Command {
	return &LoopCmd{
		name:        "loop",
		group:       "player",
		summary:     "Loop a track or the queue",
		description: "This command loops either the current track or the entire queue. Run without any arguments for the current track, or loop all/queue/list for the whole queue. Run the commands again, or loop off, to turn off looping.",
		usage:       "loop [single|all|off]",
		aliases:     []string{""},
		voiceOnly:   true,
	}
}

func (cmd *LoopCmd) GetName() string {
	return cmd.name
}

func (cmd *LoopCmd) GetGroup() string {
	return cmd.group
}

func (cmd *LoopCmd) GetSummary() string {
	return cmd.summary
}

func (cmd *LoopCmd) GetDescription() string {
	return cmd.description
}

func (cmd *LoopCmd) GetUsage() string {
	return cmd.usage
}

func (cmd *LoopCmd) GetAliases() []string {
	return cmd.aliases
}

func (cmd *LoopCmd) IsVoiceOnlyCmd() bool {
	return cmd.voiceOnly
}

func (cmd *LoopCmd) Execute(deps CommandDependencies, opts ...string) error {
	if len(opts) == 0 {
		cmd.loopWithNoOptions(deps)
		return nil
	}

	switch opts[0] {
	case "single":
		cmd.loopSingle(deps)
		return nil
	case "all":
		fallthrough
	case "list":
		fallthrough
	case "queue":
		cmd.loopQueue(deps)
		return nil
	default:
		deps.MusicPlayer.SetLoopModeOff()
		deps.Messenger.SendSimpleMessage("Looping off.")
		return nil
	}
}

func (cmd *LoopCmd) loopWithNoOptions(deps CommandDependencies) {
	if !deps.MusicPlayer.IsLoopModeOff() {
		deps.MusicPlayer.SetLoopModeOff()
		deps.Messenger.SendSimpleMessage("Looping off.")
		return
	}

	deps.MusicPlayer.SetLoopModeTrack()
	deps.Messenger.SendSimpleMessage("Looping the current track.")
}

func (cmd *LoopCmd) loopSingle(deps CommandDependencies) {
	if deps.MusicPlayer.IsLoopModeTrack() {
		deps.MusicPlayer.SetLoopModeOff()
		deps.Messenger.SendSimpleMessage("Looping off.")
		return
	}

	deps.MusicPlayer.SetLoopModeTrack()
	deps.Messenger.SendSimpleMessage("Looping the current track.")
}

func (cmd *LoopCmd) loopQueue(deps CommandDependencies) {
	if deps.MusicPlayer.IsLoopModeQueue() {
		deps.MusicPlayer.SetLoopModeOff()
		deps.Messenger.SendSimpleMessage("Looping off.")
		return
	}

	deps.MusicPlayer.SetLoopModeQueue()
	deps.Messenger.SendSimpleMessage("Looping the whole queue.")
}
