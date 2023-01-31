package command

import (
	"schoperation/schopyatch/musicplayer"
	"schoperation/schopyatch/util"
)

type LoopCmd struct {
	Name        string
	Summary     string
	Description string
	Usage       string
	Aliases     []string
}

func NewLoopCmd() Command {
	return &LoopCmd{
		Name:        "loop",
		Summary:     "Loop a track or the queue",
		Description: "This command loops either the current track or the entire queue. Run without any arguments for the current track, or loop all/queue/list for the whole queue. Run the commands again, or loop off, to turn off looping.",
		Usage:       "loop [single|all|off]",
		Aliases:     []string{""},
	}
}

func (cmd *LoopCmd) GetName() string {
	return cmd.Name
}

func (cmd *LoopCmd) GetSummary() string {
	return cmd.Summary
}

func (cmd *LoopCmd) GetDescription() string {
	return cmd.Description
}

func (cmd *LoopCmd) GetUsage() string {
	return cmd.Usage
}

func (cmd *LoopCmd) GetAliases() []string {
	return cmd.Aliases
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
		deps.MusicPlayer.LoopMode = musicplayer.LoopOff
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Looping off.")
		return nil
	}
}

func (cmd *LoopCmd) loopSingle(deps CommandDependencies) {
	if deps.MusicPlayer.LoopMode != musicplayer.LoopOff {
		deps.MusicPlayer.LoopMode = musicplayer.LoopOff
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Looping off.")
		return
	}

	deps.MusicPlayer.LoopMode = musicplayer.LoopTrack
	util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Looping the current track.")
}

func (cmd *LoopCmd) loopQueue(deps CommandDependencies) {
	if deps.MusicPlayer.LoopMode == musicplayer.LoopQueue {
		deps.MusicPlayer.LoopMode = musicplayer.LoopOff
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Looping off.")
		return
	}

	deps.MusicPlayer.LoopMode = musicplayer.LoopQueue
	util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Looping the whole queue.")
}
