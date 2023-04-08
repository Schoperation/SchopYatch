package command

import (
	"fmt"
	"log"
	"schoperation/schopyatch/util"
	"strconv"
	"strings"
)

type QueueCmd struct {
	name        string
	summary     string
	description string
	usage       string
	aliases     []string
	voiceOnly   bool
}

func NewQueueCmd() Command {
	return &QueueCmd{
		name:        "queue",
		summary:     "View the player's queue",
		description: "Use this command to view a list of tracks that will eventually be played on the bot. If there are more than 10 tracks in the queue, you can add a page parameter to see additional pages.",
		usage:       "queue [page]",
		aliases:     []string{"q", "list"},
		voiceOnly:   false,
	}
}

func (cmd *QueueCmd) GetName() string {
	return cmd.name
}

func (cmd *QueueCmd) GetSummary() string {
	return cmd.summary
}

func (cmd *QueueCmd) GetDescription() string {
	return cmd.description
}

func (cmd *QueueCmd) GetUsage() string {
	return cmd.usage
}

func (cmd *QueueCmd) GetAliases() []string {
	return cmd.aliases
}

func (cmd *QueueCmd) IsVoiceOnlyCmd() bool {
	return cmd.voiceOnly
}

func (cmd *QueueCmd) Execute(deps CommandDependencies, opts ...string) error {
	builder := strings.Builder{}

	currentTrack, err := deps.MusicPlayer.GetLoadedTrack()
	if err == nil {
		builder.WriteString(fmt.Sprintf("Now Playing:\n\t*%s* by **%s** `[%s / %s]`\n\n", currentTrack.Info.Title, currentTrack.Info.Author, deps.MusicPlayer.GetPosition().String(), currentTrack.Info.Length.String()))
	}

	if deps.MusicPlayer.IsLoopModeTrack() {
		builder.WriteString("**Looping Current Track**\n")
	} else if deps.MusicPlayer.IsLoopModeQueue() {
		builder.WriteString("**Looping Queue**\n")
	}

	if deps.MusicPlayer.IsQueueEmpty() {
		builder.WriteString("Queue is empty.\n")
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, builder.String())
		return nil
	}

	queueLen := deps.MusicPlayer.GetQueueLength()
	pages := (queueLen / 10) + 1
	if queueLen%10 == 0 {
		pages--
	}

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

	if queueLen == 1 {
		builder.WriteString(fmt.Sprintf("Total of **%d** track in the queue. `[%s]`\n", 1, deps.MusicPlayer.GetQueueDuration().String()))
	} else {
		builder.WriteString(fmt.Sprintf("Total of **%d** tracks in the queue. `[%s]`\n", queueLen, deps.MusicPlayer.GetQueueDuration().String()))
	}

	builder.WriteString(fmt.Sprintf("Page **%d** of **%d**:\n\n", pageNum, pages))

	rangeStart := (pageNum - 1) * 10
	rangeEnd := pageNum * 10
	if rangeEnd > queueLen {
		rangeEnd = queueLen
	}

	queue := deps.MusicPlayer.GetQueue()
	for i := rangeStart; i < rangeEnd; i++ {
		builder.WriteString(fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", i+1, queue[i].Info.Title, queue[i].Info.Author, queue[i].Info.Length))
	}

	util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, builder.String())
	return nil
}
