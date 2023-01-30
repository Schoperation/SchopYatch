package command

import (
	"fmt"
	"schoperation/schopyatch/util"
	"strings"
)

type NowPlayingCmd struct {
	Name        string
	Summary     string
	Description string
	Usage       string
	Aliases     []string
}

func NewNowPlayingCmd() Command {
	return &NowPlayingCmd{
		Name:        "nowplaying",
		Summary:     "Show the details of the track currently playing",
		Description: "This command shows details about the track that's currently playing on the bot.",
		Usage:       "nowplaying",
		Aliases:     []string{"np"},
	}
}

func (cmd *NowPlayingCmd) GetName() string {
	return cmd.Name
}

func (cmd *NowPlayingCmd) GetSummary() string {
	return cmd.Summary
}

func (cmd *NowPlayingCmd) GetDescription() string {
	return cmd.Description
}

func (cmd *NowPlayingCmd) GetUsage() string {
	return cmd.Usage
}

func (cmd *NowPlayingCmd) GetAliases() []string {
	return cmd.Aliases
}

func (cmd *NowPlayingCmd) Execute(deps CommandDependencies, opts ...string) error {
	currentTrack := deps.MusicPlayer.Player.PlayingTrack()
	if currentTrack == nil {
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Nothing's playing. Bruh moment...")
		return nil
	}

	currentPos := deps.MusicPlayer.Player.Position()
	trackLen := currentTrack.Info().Length

	daysLeft := trackLen.Days() - currentPos.Days()
	hoursLeft := trackLen.HoursPart() - currentPos.HoursPart()
	minutesLeft := trackLen.MinutesPart() - currentPos.MinutesPart()
	secondsLeft := trackLen.SecondsPart() - currentPos.SecondsPart()

	builder := strings.Builder{}
	if daysLeft > 1 {
		builder.WriteString(fmt.Sprintf("%d days, ", daysLeft))
	} else if daysLeft == 1 {
		builder.WriteString(fmt.Sprintf("%d day, ", daysLeft))
	}
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
	if secondsLeft > 1 {
		builder.WriteString(fmt.Sprintf("%d secs remaining.", secondsLeft))
	} else if secondsLeft == 1 {
		builder.WriteString(fmt.Sprintf("%d sec remaining.", secondsLeft))
	}

	finalStr := fmt.Sprintf("Now Playing:\n\t*%s* by **%s**\n\t%s `[%s / %s]`", currentTrack.Info().Title, currentTrack.Info().Author, builder.String(), currentPos, trackLen)
	util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, finalStr)
	return nil
}
