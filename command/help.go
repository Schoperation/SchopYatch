package command

import (
	"fmt"
	"strings"
)

type HelpCmd struct {
	name        string
	summary     string
	description string
	usage       string
	aliases     []string
	voiceOnly   bool
}

func NewHelpCmd() Command {
	return &HelpCmd{
		name:        "help",
		summary:     "Shows info about the commands",
		description: "Woah, you need a lot of help if you're asking for it twice. Have you considered middle school?",
		usage:       "help [command]",
		aliases:     []string{"h", "helpme", "thefuckisthisshit", "command", "commands"},
		voiceOnly:   false,
	}
}

func (cmd *HelpCmd) GetName() string {
	return cmd.name
}

func (cmd *HelpCmd) GetSummary() string {
	return cmd.summary
}

func (cmd *HelpCmd) GetDescription() string {
	return cmd.description
}

func (cmd *HelpCmd) GetUsage() string {
	return cmd.usage
}

func (cmd *HelpCmd) GetAliases() []string {
	return cmd.aliases
}

func (cmd *HelpCmd) IsVoiceOnlyCmd() bool {
	return cmd.voiceOnly
}

func (cmd *HelpCmd) Execute(deps CommandDependencies, opts ...string) error {
	commands := GetCommands()
	builder := strings.Builder{}

	if len(opts) > 0 {
		mappedCommands := GetCommandsAndAliasesAsMap()

		cmd, exists := mappedCommands[strings.ToLower(opts[0])]
		if !exists {
			deps.Messenger.SendSimpleMessage(fmt.Sprintf("Could not find %s. Try doing %shelp for a full list.", strings.ToLower(opts[0]), deps.Prefix))
			return nil
		}

		builder.WriteString(fmt.Sprintf("Usage: `%s%s`\n", deps.Prefix, cmd.GetUsage()))
		builder.WriteString("Aliases: `")

		for i, alias := range cmd.GetAliases() {
			if i < len(cmd.GetAliases())-1 {
				builder.WriteString(fmt.Sprintf("%s, ", alias))
			} else {
				builder.WriteString(alias)
			}
		}

		builder.WriteString(" `\n\n")

		builder.WriteString(fmt.Sprintf("%s\n", cmd.GetDescription()))

		deps.Messenger.SendSimpleMessage(builder.String())
		return nil
	}

	builder.WriteString("Hey, SchopYatch here!\n")
	builder.WriteString(fmt.Sprintf("To get started quickly, use the `%splay` command with either a search query or a URL, and I'll play it immediately in your channel.\n\n", deps.Prefix))
	builder.WriteString(fmt.Sprintf("E.g. `%splay ace attorney all pursuit themes`\nOr... `%splay https://www.youtube.com/watch?v=dv13gl0a-FA`\nAlso works with playlists and livestreams!\n\n", deps.Prefix, deps.Prefix))
	builder.WriteString(fmt.Sprintf("Use `%shelp command` to learn more about a command, e.g. `%shelp play`.\n", deps.Prefix, deps.Prefix))
	builder.WriteString(fmt.Sprintf("On those pages, usages use the following convention: `%scommand <required parameter> [optional parameter]`\n\n```", deps.Prefix))

	for i, cmd := range commands {
		builder.WriteString(fmt.Sprintf("\t%s%s ~ %s\n", deps.Prefix, cmd.GetName(), cmd.GetSummary()))

		if (i+1)%4 == 0 {
			builder.WriteString("\n")
		}
	}

	builder.WriteString("```")

	deps.Messenger.SendSimpleMessage(builder.String())
	return nil
}
