package command

var (
	helpCmd       = NewHelpCmd()
	playCmd       = NewPlayCmd()
	joinCmd       = NewJoinCmd()
	leaveCmd      = NewLeaveCmd()
	pauseCmd      = NewPauseCmd()
	resumeCmd     = NewResumeCmd()
	nowPlayingCmd = NewNowPlayingCmd()
	skipCmd       = NewSkipCmd()
	skipToCmd     = NewSkipToCmd()
	loopCmd       = NewLoopCmd()
	queueCmd      = NewQueueCmd()
	clearCmd      = NewClearCmd()
	shuffleCmd    = NewShuffleCmd()
	pingCmd       = NewPingCmd()
)

func GetCommands() []Command {
	return []Command{
		helpCmd,
		playCmd,
		joinCmd,
		leaveCmd,
		pauseCmd,
		resumeCmd,
		nowPlayingCmd,
		skipCmd,
		skipToCmd,
		loopCmd,
		queueCmd,
		clearCmd,
		shuffleCmd,
		pingCmd,
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
