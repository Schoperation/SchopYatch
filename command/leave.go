package command

import (
	"context"
	"schoperation/schopyatch/util"
)

type LeaveCmd struct {
	Name        string
	Summary     string
	Description string
	Usage       string
	Aliases     []string
}

func NewLeaveCmd() Command {
	return &LeaveCmd{
		Name:        "leave",
		Summary:     "Make the bot leave a voice channel",
		Description: "Upon running, the bot will leave the user's voice channel. Kindly.",
		Usage:       "leave",
		Aliases:     []string{"fuckoff"},
	}
}

func (cmd *LeaveCmd) GetName() string {
	return cmd.Name
}

func (cmd *LeaveCmd) GetSummary() string {
	return cmd.Summary
}

func (cmd *LeaveCmd) GetDescription() string {
	return cmd.Description
}

func (cmd *LeaveCmd) GetUsage() string {
	return cmd.Usage
}

func (cmd *LeaveCmd) GetAliases() []string {
	return cmd.Aliases
}

func (cmd *LeaveCmd) Execute(deps CommandDependencies, opts ...string) error {
	err := leaveVoiceChannel(deps)
	return err
}

func leaveVoiceChannel(deps CommandDependencies) error {
	if deps.MusicPlayer.Player.PlayingTrack() != nil {
		err := deps.MusicPlayer.Player.Stop()
		if err != nil {
			util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "For some reason I can't stop this track...")
			return err
		}
	}

	err := (*deps.Client).UpdateVoiceState(context.TODO(), *deps.Event.GuildID, nil, false, true)
	if err != nil {
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Cannot leave the channel... wait what?")
		return err
	}

	deps.MusicPlayer.Queue.Clear()
	deps.MusicPlayer.ClearSearchResults()
	deps.MusicPlayer.GotDisconnected = true
	return err
}
