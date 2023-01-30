package command

import (
	"fmt"
	"schoperation/schopyatch/util"
)

type SkipCmd struct {
	Name        string
	Summary     string
	Description string
	Usage       string
	Aliases     []string
}

func NewSkipCmd() Command {
	return &SkipCmd{
		Name:        "skip",
		Summary:     "Skip the current track",
		Description: "This command skips the track that's currently playing on the bot. If the queue has tracks, it'll play the next one in line. Otherwise, it'll go radio silent...",
		Usage:       "skip",
		Aliases:     []string{"s", "next"},
	}
}

func (cmd *SkipCmd) GetName() string {
	return cmd.Name
}

func (cmd *SkipCmd) GetSummary() string {
	return cmd.Summary
}

func (cmd *SkipCmd) GetDescription() string {
	return cmd.Description
}

func (cmd *SkipCmd) GetUsage() string {
	return cmd.Usage
}

func (cmd *SkipCmd) GetAliases() []string {
	return cmd.Aliases
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

	nextTrack := deps.MusicPlayer.Queue.Dequeue()
	err := deps.MusicPlayer.Player.Play(*nextTrack)
	if err != nil {
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "For some reason I can't play this... might be some dumb age restriction?")
		return err
	}

	util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, fmt.Sprintf("Now playing *%s* by **%s**.", (*nextTrack).Info().Title, (*nextTrack).Info().Author))
	return nil
}
