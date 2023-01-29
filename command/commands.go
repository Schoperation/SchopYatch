package command

var (
	helpCmd = NewHelpCmd()
	pingCmd = NewPingCmd()
	playCmd = NewPlayCmd()
)

func GetCommands() []Command {
	return []Command{
		helpCmd,
		pingCmd,
		playCmd,
	}
}

func GetCommandsAndAliasesAsMap() map[string]Command {
	commands := GetCommands()
	var mappedCommands = make(map[string]Command)

	for _, command := range commands {
		mappedCommands[command.GetName()] = command

		if len(command.GetAliases()) > 0 {
			for _, alias := range command.GetAliases() {
				mappedCommands[alias] = command
			}
		}
	}

	return mappedCommands
}
