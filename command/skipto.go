package command

import (
	"fmt"
	"schoperation/schopyatch/msg"
	"strconv"
)

type SkipToCmd struct {
	name        string
	group       string
	summary     string
	description string
	usage       string
	aliases     []string
	voiceOnly   bool
}

func NewSkipToCmd() Command {
	return &SkipToCmd{
		name:        "skipto",
		group:       "player",
		summary:     "Skip multiple tracks to a spot in the queue",
		description: "This command skips not only the currently playing track, but any tracks in the queue that are before the specified position number. To see the position numbers, run the queue command.",
		usage:       "skipto <position>",
		aliases:     []string{"st", "sto", "nextto"},
		voiceOnly:   true,
	}
}

func (cmd *SkipToCmd) GetName() string {
	return cmd.name
}

func (cmd *SkipToCmd) GetGroup() string {
	return cmd.group
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

func (cmd *SkipToCmd) IsVoiceOnlyCmd() bool {
	return cmd.voiceOnly
}

func (cmd *SkipToCmd) Execute(deps CommandDependencies, opts ...string) error {
	if len(opts) == 0 {
		deps.Messenger.SendSimpleMessage("No position specified. Please specify a position in the queue. E.g. `skipto 5` to go to the 5th song in the queue.")
		return nil
	}

	num, err := strconv.Atoi(opts[0])
	if err != nil {
		deps.Messenger.SendSimpleMessage("Woah hey now, that ain't a number...")
		return nil
	}

	deps.Messenger.SendSimpleMessage(fmt.Sprintf("Skipping to #%d in the queue...", num))

	playingTrack, err := deps.MusicPlayer.SkipTo(num - 1)
	if err != nil {
		if msg.IsErrorMessage(err, msg.NoLoadedTrack) {
			deps.Messenger.SendSimpleMessage("Nothing to skip. Have a great evening.")
			return nil
		}
		if msg.IsErrorMessage(err, msg.IndexOutOfBounds) {
			deps.Messenger.SendSimpleMessage(fmt.Sprintf("Out of bounds. Please use a number between 1 and %d", deps.MusicPlayer.GetQueueLength()))
			return nil
		}

		return err
	}

	if playingTrack == nil {
		deps.Messenger.SendSimpleMessage("All is now quiet on the SchopYatch front.")
		return nil
	}

	deps.Messenger.SendSimpleMessage(fmt.Sprintf("Now playing *%s* by **%s**.", playingTrack.Title(), playingTrack.Author()))
	return nil
}
