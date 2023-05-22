package command

import (
	"fmt"
	"strings"
)

type AboutCmd struct {
	name        string
	summary     string
	description string
	usage       string
	aliases     []string
	voiceOnly   bool
}

func NewAboutCmd() Command {
	return &AboutCmd{
		name:        "about",
		summary:     "Show bot information",
		description: "Displays a variety of information about the bot, including version, github link, and more.",
		usage:       "about",
		voiceOnly:   false,
	}
}

func (cmd *AboutCmd) GetName() string {
	return cmd.name
}

func (cmd *AboutCmd) GetSummary() string {
	return cmd.summary
}

func (cmd *AboutCmd) GetDescription() string {
	return cmd.description
}

func (cmd *AboutCmd) GetUsage() string {
	return cmd.usage
}

func (cmd *AboutCmd) GetAliases() []string {
	return cmd.aliases
}

func (cmd *AboutCmd) IsVoiceOnlyCmd() bool {
	return cmd.voiceOnly
}

func (cmd *AboutCmd) Execute(deps CommandDependencies, opts ...string) error {
	builder := strings.Builder{}

	builder.WriteString("```")
	builder.WriteString(fmt.Sprintf("SchopYatch v%s\n\n", deps.Version))
	builder.WriteString("Coded by Schoperation: 		   https://github.com/Schoperation/SchopYatch\n")
	builder.WriteString("Lavalink Client by Freya Arbjerg: https://github.com/freyacodes/Lavalink-Client\n")
	builder.WriteString("Libraries written by the DisGoOrg:\n")
	builder.WriteString("\tDisGo:     https://github.com/DisgoOrg/disgo\n")
	builder.WriteString("\tDisGoLink: https://github.com/disgoorg/disgolink\n")
	builder.WriteString("```")

	deps.Messenger.SendSimpleMessage(builder.String())
	return nil
}
