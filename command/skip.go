package command

import (
	"fmt"
	"schoperation/schopyatch/musicplayer"
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
	if deps.MusicPlayer.Player.PlayingTrack() == nil {
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Nothing to skip. Have a great evening.")
		return nil
	}

	util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Skipping...")

	if deps.MusicPlayer.Queue.IsEmpty() {
		err := deps.MusicPlayer.Player.Stop()
		if err != nil {
			util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "For some reason I can't stop this track...")
			return err
		}

		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "All is now quiet on the SchopYatch front.")
		return nil
	}

	if deps.MusicPlayer.LoopMode == musicplayer.LoopQueue {
		deps.MusicPlayer.Queue.Enqueue(deps.MusicPlayer.Player.PlayingTrack())
	}

	nextTrack := deps.MusicPlayer.Queue.Dequeue()
	err := deps.MusicPlayer.Player.Play(*nextTrack)
	if err != nil {
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "For some reason I can't play this... might be some dumb age restriction?")
		return err
	}

	util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, fmt.Sprintf("Now playing *%s* by **%s**.", (*nextTrack).Info().Title, (*nextTrack).Info().Author))
	return nil
}
