package command

import (
	"fmt"
	"schoperation/schopyatch/util"
	"strconv"
	"strings"

	"github.com/disgoorg/disgolink/lavalink"
)

type SeekCmd struct {
	name        string
	summary     string
	description string
	usage       string
	aliases     []string
	voiceOnly   bool
}

func NewSeekCmd() Command {
	return &SeekCmd{
		name:        "seek",
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
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, fmt.Sprintf("Need to specify a time to seek to. `%sseek hh:mm:ss`", deps.Prefix))
		return nil
	}

	splitString := strings.Split(opts[0], ":")

	var times []int
	for _, str := range splitString {
		time, err := strconv.Atoi(str)
		if err != nil {
			util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, fmt.Sprintf("Detected something that ain't a number... cmon man... `%sseek hh:mm:ss`", deps.Prefix))
			return err
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
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "OBJECTION! Your extra times contradict the ISO standard!")
		return nil
	}

	if duration >= deps.MusicPlayer.Player.PlayingTrack().Info().Length {
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, fmt.Sprintf("The time you specified is longer than the length of `%s`. Try using your fingers.", deps.MusicPlayer.Player.PlayingTrack().Info().Length))
		return nil
	}

	err := deps.MusicPlayer.Player.Seek(duration)
	if err != nil {
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Couldn't seek out some sikh moves. Maybe the times are off?")
		return err
	}

	util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, fmt.Sprintf("Seeked to `%s`", duration.String()))
	return nil
}
