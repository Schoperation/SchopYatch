package command

import (
	"fmt"
	"schoperation/schopyatch/musicplayer"
	"schoperation/schopyatch/util"
	"strconv"
)

type SkipToCmd struct {
	name        string
	summary     string
	description string
	usage       string
	aliases     []string
}

func NewSkipToCmd() Command {
	return &SkipToCmd{
		name:        "skipto",
		summary:     "Skip multiple tracks to a spot in the queue",
		description: "This command skips not only the currently playing track, but any tracks in the queue that are before the specified position number. To see the position numbers, run the queue command.",
		usage:       "skipto <position>",
		aliases:     []string{"st", "sto", "nextto"},
	}
}

func (cmd *SkipToCmd) GetName() string {
	return cmd.name
}

func (cmd *SkipToCmd) GetSummary() string {
	return cmd.summary
}

func (cmd *SkipToCmd) GetDescription() string {
	return cmd.description
}

func (cmd *SkipToCmd) GetUsage() string {
	return cmd.usage
}

func (cmd *SkipToCmd) GetAliases() []string {
	return cmd.aliases
}

func (cmd *SkipToCmd) Execute(deps CommandDependencies, opts ...string) error {
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

	if len(opts) == 0 {
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "No position specified. Please specify a position in the queue. E.g. `skipto 5` to go to the 5th song in the queue.")
		return nil
	}

	if deps.MusicPlayer.LoopMode == musicplayer.LoopQueue {
		deps.MusicPlayer.Queue.Enqueue(deps.MusicPlayer.Player.PlayingTrack())
	}

	num, err := strconv.Atoi(opts[0])
	if err != nil {
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Woah hey now, that ain't a number...")
		return err
	}

	if num > deps.MusicPlayer.Queue.Length() {
		num = deps.MusicPlayer.Queue.Length()
	}

	for i := 0; i < num-1; i++ {
		track := deps.MusicPlayer.Queue.Dequeue()

		if deps.MusicPlayer.LoopMode == musicplayer.LoopQueue {
			deps.MusicPlayer.Queue.Enqueue(*track)
		}
	}

	nextTrack := deps.MusicPlayer.Queue.Dequeue()
	err = deps.MusicPlayer.Player.Play(*nextTrack)
	if err != nil {
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "For some reason I can't play this... might be some dumb age restriction?")
		return err
	}

	util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, fmt.Sprintf("Now playing *%s* by **%s**.", (*nextTrack).Info().Title, (*nextTrack).Info().Author))
	return nil
}
