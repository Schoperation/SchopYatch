package command

import (
	"fmt"
	"schoperation/schopyatch/util"
)

type SkipCmd struct {
	name        string
	summary     string
	description string
	usage       string
	aliases     []string
	voiceOnly   bool
}

func NewSkipCmd() Command {
	return &SkipCmd{
		name:        "skip",
		summary:     "Skip the current track",
		description: "This command skips the track that's currently playing on the bot. If the queue has tracks, it'll play the next one in line. Otherwise, it'll go radio silent...",
		usage:       "skip",
		aliases:     []string{"s", "next"},
		voiceOnly:   true,
	}
}

func (cmd *SkipCmd) GetName() string {
	return cmd.name
}

func (cmd *SkipCmd) GetSummary() string {
	return cmd.summary
}

func (cmd *SkipCmd) GetDescription() string {
	return cmd.description
}

func (cmd *SkipCmd) GetUsage() string {
	return cmd.usage
}

func (cmd *SkipCmd) GetAliases() []string {
	return cmd.aliases
}

func (cmd *SkipCmd) IsVoiceOnlyCmd() bool {
	return cmd.voiceOnly
}

func (cmd *SkipCmd) Execute(deps CommandDependencies, opts ...string) error {
	util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Skipping...")

	playingTrack, err := deps.MusicPlayer.Skip()
	if err != nil && util.IsErrorMessage(err, util.NoLoadedTrack) {
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Nothing to skip. Have a great evening.")
		return nil
	} else if err != nil {
		return err
	}

	if playingTrack == nil {
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "All is now quiet on the SchopYatch front.")
		return nil
	}

	util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, fmt.Sprintf("Now playing *%s* by **%s**.", playingTrack.Info.Title, playingTrack.Info.Author))
	return nil
}
