package command

import (
	"fmt"
	"log"
	"schoperation/schopyatch/util"
	"strconv"
	"strings"
)

type QueueCmd struct {
	Name        string
	Summary     string
	Description string
	Usage       string
	Aliases     []string
}

func NewQueueCmd() Command {
	return &QueueCmd{
		Name:        "queue",
		Summary:     "View the player's queue",
		Description: "Use this command to view a list of tracks that will eventually be played on the bot. If there are more than 10 tracks in the queue, you can add a page parameter to see additional pages.",
		Usage:       "queue [page]",
		Aliases:     []string{"q", "list"},
	}
}

func (cmd *QueueCmd) GetName() string {
	return cmd.Name
}

func (cmd *QueueCmd) GetSummary() string {
	return cmd.Summary
}

func (cmd *QueueCmd) GetDescription() string {
	return cmd.Description
}

func (cmd *QueueCmd) GetUsage() string {
	return cmd.Usage
}

func (cmd *QueueCmd) GetAliases() []string {
	return cmd.Aliases
}

func (cmd *QueueCmd) Execute(deps CommandDependencies, opts ...string) error {
	builder := strings.Builder{}

	if deps.MusicPlayer.Player.PlayingTrack() != nil {
		currentTrack := deps.MusicPlayer.Player.PlayingTrack()
		builder.WriteString(fmt.Sprintf("Now playing: *%s* by **%s** `[%s / %s]`\n\n", currentTrack.Info().Title, currentTrack.Info().Author, deps.MusicPlayer.Player.Position().String(), currentTrack.Info().Length.String()))
	}

	if deps.MusicPlayer.Queue.IsEmpty() {
		builder.WriteString("Queue is empty.\n")
		err := util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, builder.String())
		return err
	}

	pages := (deps.MusicPlayer.Queue.Length() / 10) + 1
	pageNum := 1
	if len(opts) > 0 {
		num, err := strconv.Atoi(opts[0])
		if err != nil {
			log.Printf("Couldn't read number in queue command, ignoring...")
		} else if num > pages {
			log.Printf("Some guy tried to go out of bounds in queue command, ignoring...")
		} else {
			pageNum = num
		}
	}
	builder.WriteString(fmt.Sprintf("Page **%d** of **%d** of the queue:\n\n", pageNum, pages))

	rangeStart := (pageNum - 1) * 10
	rangeEnd := pageNum * 10
	if rangeEnd > deps.MusicPlayer.Queue.Length() {
		rangeEnd = deps.MusicPlayer.Queue.Length()
	}

	queue := deps.MusicPlayer.Queue.PeekList()
	for i := rangeStart; i < rangeEnd; i++ {
		builder.WriteString(fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", i+1, queue[i].Info().Title, queue[i].Info().Author, queue[i].Info().Length))
	}

	err := util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, builder.String())
	return err
}
