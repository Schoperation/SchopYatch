package command

import (
	"context"
	"schoperation/schopyatch/util"
)

type LeaveCmd struct {
	name        string
	summary     string
	description string
	usage       string
	aliases     []string
}

func NewLeaveCmd() Command {
	return &LeaveCmd{
		name:        "leave",
		summary:     "Make the bot leave a voice channel",
		description: "Upon running, the bot will leave the user's voice channel. Kindly.",
		usage:       "leave",
		aliases:     []string{"fuckoff"},
	}
}

func (cmd *LeaveCmd) GetName() string {
	return cmd.name
}

func (cmd *LeaveCmd) GetSummary() string {
	return cmd.summary
}

func (cmd *LeaveCmd) GetDescription() string {
	return cmd.description
}

func (cmd *LeaveCmd) GetUsage() string {
	return cmd.usage
}

func (cmd *LeaveCmd) GetAliases() []string {
	return cmd.aliases
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
