package command

import (
	"fmt"
	"schoperation/schopyatch/msg"
	"strconv"
	"strings"

	"github.com/disgoorg/disgolink/v2/lavalink"
)

type SeekCmd struct {
	name        string
	group       string
	summary     string
	description string
	usage       string
	aliases     []string
	voiceOnly   bool
}

func NewSeekCmd() Command {
	return &SeekCmd{
		name:        "seek",
		group:       "player",
		summary:     "Seek to a position in the current track",
		description: "This command allows you to seek to a specific time within the currently playing track. Ex. `seek 30` goes to 30 seconds, `seek 1:30` goes to 1 minute and 30 seconds, and `seek 02:01:30` goes to 2 hours, 1 minute, and 30 seconds. Leading zeros (01) are optional.",
		usage:       "seek <hh:mm:ss>",
		aliases:     []string{""},
		voiceOnly:   true,
	}
}

func (cmd *SeekCmd) GetName() string {
	return cmd.name
}

func (cmd *SeekCmd) GetGroup() string {
	return cmd.group
}

func (cmd *SeekCmd) GetSummary() string {
	return cmd.summary
}

func (cmd *SeekCmd) GetDescription() string {
	return cmd.description
}

func (cmd *SeekCmd) GetUsage() string {
	return cmd.usage
}

func (cmd *SeekCmd) GetAliases() []string {
	return cmd.aliases
}

func (cmd *SeekCmd) IsVoiceOnlyCmd() bool {
	return cmd.voiceOnly
}

func (cmd *SeekCmd) Execute(deps CommandDependencies, opts ...string) error {
	if len(opts) == 0 {
		deps.Messenger.SendSimpleMessage(fmt.Sprintf("Need to specify a time to seek to. `%sseek hh:mm:ss`", deps.Prefix))
		return nil
	}

	splitString := strings.Split(opts[0], ":")

	var times []int
	for _, str := range splitString {
		time, err := strconv.Atoi(str)
		if err != nil {
			deps.Messenger.SendSimpleMessage(fmt.Sprintf("Detected something that ain't a number... cmon man... `%sseek hh:mm:ss`", deps.Prefix))
			return nil
		}

		times = append(times, time)
	}

	var duration lavalink.Duration
	switch len(times) {
	case 1:
		duration = lavalink.Duration(times[0] * int(lavalink.Second))
	case 2:
		duration = lavalink.Duration(times[0]*int(lavalink.Minute) + times[1]*int(lavalink.Second))
	case 3:
		duration = lavalink.Duration(times[0]*int(lavalink.Hour) + times[1]*int(lavalink.Minute) + times[2]*int(lavalink.Second))
	default:
		deps.Messenger.SendSimpleMessage("OBJECTION! Your extra times contradict the ISO standard!")
		return nil
	}

	_, err := deps.MusicPlayer.Seek(duration)
	if err != nil {
		if msg.IsErrorMessage(err, msg.NoLoadedTrack) {
			deps.Messenger.SendSimpleMessage("Can't seek through nothing... not a sikh move if you ask me...")
			return nil
		}
		if msg.IsErrorMessage(err, msg.NegativeDuration) {
			deps.Messenger.SendSimpleMessage("Can't seek back in time. Try adopting a more positive outlook.")
			return nil
		}
		if msg.IsErrorMessage(err, msg.IndexOutOfBounds) {
			track, err := deps.MusicPlayer.GetLoadedTrack()
			if err != nil {
				return err
			}

			deps.Messenger.SendSimpleMessage(fmt.Sprintf("Can't seek beyond the track length. It's %s long.", track.Info.Length))
			return nil
		}

		return err
	}

	deps.Messenger.SendSimpleMessage(fmt.Sprintf("Seeked to `%s`", duration.String()))
	return nil
}
