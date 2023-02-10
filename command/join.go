package command

import (
	"context"
	"errors"
	"schoperation/schopyatch/util"
)

type JoinCmd struct {
	name        string
	summary     string
	description string
	usage       string
	aliases     []string
	voiceOnly   bool
}

func NewJoinCmd() Command {
	return &JoinCmd{
		name:        "join",
		summary:     "Make the bot join a voice channel",
		description: "Upon running, the bot will join the user's voice channel. It will error out if either the user isn't in a voice channel, or if the bot doesn't have permission to join.",
		usage:       "join",
		aliases:     []string{"j", "summon"},
		voiceOnly:   true,
	}
}

func (cmd *JoinCmd) GetName() string {
	return cmd.name
}

func (cmd *JoinCmd) GetSummary() string {
	return cmd.summary
}

func (cmd *JoinCmd) GetDescription() string {
	return cmd.description
}

func (cmd *JoinCmd) GetUsage() string {
	return cmd.usage
}

func (cmd *JoinCmd) GetAliases() []string {
	return cmd.aliases
}

func (cmd *JoinCmd) IsVoiceOnlyCmd() bool {
	return cmd.voiceOnly
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

	if deps.MusicPlayer.GotDisconnected {
		deps.MusicPlayer.RecreatePlayer(*deps.Lavalink)
		deps.MusicPlayer.GotDisconnected = false
	}

	return err
}
