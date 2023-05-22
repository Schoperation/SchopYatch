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
	seekCmd       = NewSeekCmd()
	loopCmd       = NewLoopCmd()
	queueCmd      = NewQueueCmd()
	clearCmd      = NewClearCmd()
	removeCmd     = NewRemoveCmd()
	shuffleCmd    = NewShuffleCmd()
	pingCmd       = NewPingCmd()
	aboutCmd      = NewAboutCmd()
)

func GetCommands() []Command {
	return []Command{
		helpCmd,
		aboutCmd,
		pingCmd,
		joinCmd,
		leaveCmd,
		playCmd,
		pauseCmd,
		resumeCmd,
		nowPlayingCmd,
		skipCmd,
		skipToCmd,
		seekCmd,
		loopCmd,
		queueCmd,
		clearCmd,
		removeCmd,
		shuffleCmd,
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
