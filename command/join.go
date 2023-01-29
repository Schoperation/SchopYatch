package command

import (
	"context"
	"errors"
	"schoperation/schopyatch/util"
)

type JoinCmd struct {
	Name        string
	Summary     string
	Description string
	Usage       string
	Aliases     []string
}

func NewJoinCmd() Command {
	return &JoinCmd{
		Name:        "join",
		Summary:     "Make the bot join a voice channel",
		Description: "Upon running, the bot will join the user's voice channel. It will error out if either the user isn't in a voice channel, or if the bot doesn't have permission to join.",
		Usage:       "join",
		Aliases:     []string{"j"},
	}
}

func (cmd *JoinCmd) GetName() string {
	return cmd.Name
}

func (cmd *JoinCmd) GetSummary() string {
	return cmd.Summary
}

func (cmd *JoinCmd) GetDescription() string {
	return cmd.Description
}

func (cmd *JoinCmd) GetUsage() string {
	return cmd.Usage
}

func (cmd *JoinCmd) GetAliases() []string {
	return cmd.Aliases
}

func (cmd *JoinCmd) Execute(deps CommandDependencies, opts ...string) error {
	err := joinVoiceChannel(deps)
	return err
}

func joinVoiceChannel(deps CommandDependencies) error {
	voiceState, exists := (*deps.Client).Caches().VoiceState(*deps.Event.GuildID, deps.Event.Message.Author.ID)
	if !exists {
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Dude you're not in a voice channel... get in one I can see!")
		return errors.New("could not find voice state")
	}

	err := (*deps.Client).UpdateVoiceState(context.TODO(), *deps.Event.GuildID, voiceState.ChannelID, false, true)
	if err != nil {
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Cannot connect to your channel... do I have permission?")
	}

	return err
}
