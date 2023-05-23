package command

var (
	aboutCmd      = NewAboutCmd()
	clearCmd      = NewClearCmd()
	helpCmd       = NewHelpCmd()
	joinCmd       = NewJoinCmd()
	leaveCmd      = NewLeaveCmd()
	loopCmd       = NewLoopCmd()
	nowPlayingCmd = NewNowPlayingCmd()
	pauseCmd      = NewPauseCmd()
	pingCmd       = NewPingCmd()
	playCmd       = NewPlayCmd()
	queueCmd      = NewQueueCmd()
	removeCmd     = NewRemoveCmd()
	resumeCmd     = NewResumeCmd()
	seekCmd       = NewSeekCmd()
	shuffleCmd    = NewShuffleCmd()
	skipCmd       = NewSkipCmd()
	skipToCmd     = NewSkipToCmd()
)

func GetCommands() []Command {
	return []Command{
		aboutCmd,
		clearCmd,
		helpCmd,
		joinCmd,
		leaveCmd,
		loopCmd,
		nowPlayingCmd,
		pauseCmd,
		pingCmd,
		playCmd,
		queueCmd,
		removeCmd,
		resumeCmd,
		seekCmd,
		shuffleCmd,
		skipCmd,
		skipToCmd,
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
