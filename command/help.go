package command

import (
	"fmt"
	"schoperation/schopyatch/util"
	"strings"
)

type HelpCmd struct {
	Name        string
	Summary     string
	Description string
	Usage       string
	Aliases     []string
}

func NewHelpCmd() Command {
	return &HelpCmd{
		Name:        "help",
		Summary:     "Shows info about the commands",
		Description: "Woah, you need a lot of help if you're asking for it twice. Have you considered middle school?",
		Usage:       "help [command]",
		Aliases:     []string{"h", "helpme", "thefuckisthisshit"},
	}
}

func (cmd *HelpCmd) GetName() string {
	return cmd.Name
}

func (cmd *HelpCmd) GetSummary() string {
	return cmd.Summary
}

func (cmd *HelpCmd) GetDescription() string {
	return cmd.Description
}

func (cmd *HelpCmd) GetUsage() string {
	return cmd.Usage
}

func (cmd *HelpCmd) GetAliases() []string {
	return cmd.Aliases
}

func (cmd *HelpCmd) Execute(deps CommandDependencies, opts ...string) error {
	commands := GetCommands()
	builder := strings.Builder{}

	if len(opts) > 0 {
		mappedCommands := GetCommandsAndAliasesAsMap()

		cmd, exists := mappedCommands[strings.ToLower(opts[0])]
		if !exists {
			err := util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, fmt.Sprintf("Could not find %s. Try doing %shelp for a full list.", strings.ToLower(opts[0]), deps.Prefix))

			return err
		}

		builder.WriteString(fmt.Sprintf("Usage: `%s%s`\n", deps.Prefix, cmd.GetUsage()))
		builder.WriteString("Aliases: `")

		for i, alias := range cmd.GetAliases() {
			if i < len(cmd.GetAliases())-1 {
				builder.WriteString(fmt.Sprintf("%s, ", alias))
			} else {
				builder.WriteString(fmt.Sprintf("%s`\n\n", alias))
			}
		}

		builder.WriteString(fmt.Sprintf("%s\n", cmd.GetDescription()))

		err := util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, builder.String())
		return err
	}

	builder.WriteString("SchopYatch is designed similarly to the old FredYatch, albeit with a few missing features and different command names.\n")
	builder.WriteString("For a quickstart, use the play command with either a url or a search query, and the bot should play music immediately in your channel.\n\n")
	builder.WriteString(fmt.Sprintf("Full list of commands is below. To see more info about a command, you can use `%shelp command`.\n", deps.Prefix))
	builder.WriteString(fmt.Sprintf("Usages look like this: `%scommand <required parameter> [optional parameter]`\n\n```", deps.Prefix))

	for _, cmd := range commands {
		builder.WriteString(fmt.Sprintf("\t%s%s ~ %s\n", deps.Prefix, cmd.GetName(), cmd.GetSummary()))
	}

	builder.WriteString("```")

	err := util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, builder.String())
	return err
}
