package command

import (
	"fmt"
	"net/url"
	"schoperation/schopyatch/enum"
	"schoperation/schopyatch/util"
	"strconv"
	"strings"

	"github.com/disgoorg/disgolink/v2/lavalink"
)

type PlayCmd struct {
	name        string
	summary     string
	description string
	usage       string
	aliases     []string
	voiceOnly   bool
}

func NewPlayCmd() Command {
	return &PlayCmd{
		name:        "play",
		summary:     "Play a track or playlist",
		description: "Plays a track on the bot",
		usage:       "play <required> [optional]",
		aliases:     []string{"p", "resume"},
		voiceOnly:   true,
	}
}

func (cmd *PlayCmd) GetName() string {
	return cmd.name
}

func (cmd *PlayCmd) GetSummary() string {
	return cmd.summary
}

func (cmd *PlayCmd) GetDescription() string {
	return cmd.description
}

func (cmd *PlayCmd) GetUsage() string {
	return cmd.usage
}

func (cmd *PlayCmd) GetAliases() []string {
	return cmd.aliases
}

func (cmd *PlayCmd) IsVoiceOnlyCmd() bool {
	return cmd.voiceOnly
}

func (cmd *PlayCmd) Execute(deps CommandDependencies, opts ...string) error {
	if len(opts) == 0 {
		status, err := deps.MusicPlayer.Unpause()
		if err != nil {
			util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Bruh where's your song??")
		} else if status == enum.StatusAlreadyUnpaused {
			util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Already playing.")
		}

		return nil
	}

	// Selecting a search result
	num, err := strconv.Atoi(opts[0])
	if err == nil {
		if num < 1 || num > deps.MusicPlayer.GetSearchResultsLength() {
			util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, fmt.Sprintf("Selected thin air. Try a number between 1 and %d.", deps.MusicPlayer.GetSearchResultsLength()))
			return nil
		}

		track := deps.MusicPlayer.GetSearchResult(num - 1)
		status, err := deps.MusicPlayer.Load(*track)
		if err != nil {
			return err
		}

		if status == enum.StatusQueued {
			util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, fmt.Sprintf("Queued *%s* by **%s**.", track.Info.Title, track.Info.Author))
			return nil
		}

		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, fmt.Sprintf("Now playing *%s* by **%s**.", track.Info.Title, track.Info.Author))
		return nil
	}

	// If we have a search query, then we'll need to put the string opts back together into one.
	song := opts[0]
	_, err = url.ParseRequestURI(song)
	if err != nil {
		song = fmt.Sprintf("%s:%s", lavalink.SearchTypeYoutube, strings.Join(opts, " "))
	}

	err = deps.MusicPlayer.JoinVoiceChannel(deps.Client, deps.Event.Message.Author.ID)
	if err != nil {
		if util.IsErrorMessage(err, util.VoiceStateNotFound) {
			util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Dude you're not in a voice channel... get in one I can see!")
			return nil
		}

		return err
	}

	status, track, tracksQueued, err := deps.MusicPlayer.ProcessQuery(song)
	if err != nil {
		if util.IsErrorMessage(err, util.NoResultsFound) {
			util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "No results. Try some other keywords? Such as OFFICIAL, FEATURING, ft., THE TRUTH ABOUT, IS A FRAUD, or CHARLIE")
			return nil
		}

		return err
	}

	switch status {
	case enum.StatusQueued:
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, fmt.Sprintf("Queued *%s* by **%s**.", track.Info.Title, track.Info.Author))
	case enum.StatusQueuedList:
		if tracksQueued == 0 {
			util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Queued nothing. What the...?")
		} else if tracksQueued == 1 {
			util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Queued **1** additional track. Just a one hit wonder, huh?")
		} else {
			util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, fmt.Sprintf("Queued **%d** additional tracks.", tracksQueued))
		}
	case enum.StatusPlayingAndQueuedList:
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, fmt.Sprintf("Now playing *%s* by **%s**.", track.Info.Title, track.Info.Author))
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Queued **1** additional track.")
	case enum.StatusSearchSuccess:
		builder := strings.Builder{}
		builder.WriteString("Search Results:\n\n")
		for i, result := range deps.MusicPlayer.GetSearchResults() {
			builder.WriteString(fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", i+1, result.Info.Title, result.Info.Author, result.Info.Length))
			builder.WriteString(fmt.Sprintf("\nUse `%splay n` to pick a track to play.", deps.Prefix))
		}
	default:
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, fmt.Sprintf("Now playing *%s* by **%s**.", track.Info.Title, track.Info.Author))
	}

	return nil
}
