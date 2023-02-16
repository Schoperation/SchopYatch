package command

import (
	"fmt"
	"schoperation/schopyatch/musicplayer"
	"schoperation/schopyatch/util"
	"strings"
)

type NowPlayingCmd struct {
	name        string
	summary     string
	description string
	usage       string
	aliases     []string
	voiceOnly   bool
}

func NewNowPlayingCmd() Command {
	return &NowPlayingCmd{
		name:        "nowplaying",
		summary:     "Show the details of the track currently playing",
		description: "This command shows details about the track that's currently playing on the bot.",
		usage:       "nowplaying",
		aliases:     []string{"np"},
		voiceOnly:   false,
	}
}

func (cmd *NowPlayingCmd) GetName() string {
	return cmd.name
}

func (cmd *NowPlayingCmd) GetSummary() string {
	return cmd.summary
}

func (cmd *NowPlayingCmd) GetDescription() string {
	return cmd.description
}

func (cmd *NowPlayingCmd) GetUsage() string {
	return cmd.usage
}

func (cmd *NowPlayingCmd) GetAliases() []string {
	return cmd.aliases
}

func (cmd *NowPlayingCmd) IsVoiceOnlyCmd() bool {
	return cmd.voiceOnly
}

func (cmd *NowPlayingCmd) Execute(deps CommandDependencies, opts ...string) error {
	currentTrack := deps.MusicPlayer.Player.PlayingTrack()
	if currentTrack == nil {
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Nothing's playing. Bruh moment...")
		return nil
	}

	currentPos := deps.MusicPlayer.Player.Position()
	trackLen := currentTrack.Info().Length
	timeLeft := trackLen.Seconds() - currentPos.Seconds()

	hoursLeft := timeLeft / 3600
	timeLeft %= 3600
	minutesLeft := timeLeft / 60
	timeLeft %= 60
	secondsLeft := timeLeft

	builder := strings.Builder{}
	if hoursLeft > 1 {
		builder.WriteString(fmt.Sprintf("%d hrs, ", hoursLeft))
	} else if hoursLeft == 1 {
		builder.WriteString(fmt.Sprintf("%d hr, ", hoursLeft))
	}
	if minutesLeft > 1 {
		builder.WriteString(fmt.Sprintf("%d mins, ", minutesLeft))
	} else if minutesLeft == 1 {
		builder.WriteString(fmt.Sprintf("%d min, ", minutesLeft))
	}
	if secondsLeft > 1 || secondsLeft == 0 {
		builder.WriteString(fmt.Sprintf("%d secs remaining.", secondsLeft))
	} else if secondsLeft == 1 {
		builder.WriteString(fmt.Sprintf("%d sec remaining.", secondsLeft))
	}

	loopStr := ""
	if deps.MusicPlayer.LoopMode == musicplayer.LoopTrack {
		loopStr = "**Looping Current Track**\n"
	} else if deps.MusicPlayer.LoopMode == musicplayer.LoopQueue {
		loopStr = "**Looping Queue**\n"
	}

	finalStr := fmt.Sprintf("Now Playing:\n\t*%s* by **%s**\n\t%s `[%s / %s]`\n\t%s", currentTrack.Info().Title, currentTrack.Info().Author, builder.String(), currentPos, trackLen, loopStr)
	util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, finalStr)
	return nil
}
