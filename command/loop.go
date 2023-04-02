package command

import (
	"schoperation/schopyatch/music_player"
	"schoperation/schopyatch/util"
)

type LoopCmd struct {
	name        string
	summary     string
	description string
	usage       string
	aliases     []string
	voiceOnly   bool
}

func NewLoopCmd() Command {
	return &LoopCmd{
		name:        "loop",
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
		cmd.loopSingle(deps)
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
		deps.MusicPlayer.LoopMode = music_player.LoopOff
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Looping off.")
		return nil
	}
}

func (cmd *LoopCmd) loopSingle(deps CommandDependencies) {
	if deps.MusicPlayer.LoopMode != music_player.LoopOff {
		deps.MusicPlayer.LoopMode = music_player.LoopOff
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Looping off.")
		return
	}

	deps.MusicPlayer.LoopMode = music_player.LoopTrack
	util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Looping the current track.")
}

func (cmd *LoopCmd) loopQueue(deps CommandDependencies) {
	if deps.MusicPlayer.LoopMode == music_player.LoopQueue {
		deps.MusicPlayer.LoopMode = music_player.LoopOff
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Looping off.")
		return
	}

	deps.MusicPlayer.LoopMode = music_player.LoopQueue
	util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Looping the whole queue.")
}
