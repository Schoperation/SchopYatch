package command

import (
	"fmt"
	"net/url"
	"schoperation/schopyatch/enum"
	"schoperation/schopyatch/msg"
	"strconv"
	"strings"

	"github.com/disgoorg/disgolink/v3/lavalink"
)

type PlayCmd struct {
	name        string
	group       string
	summary     string
	description string
	usage       string
	aliases     []string
	voiceOnly   bool
}

func NewPlayCmd() Command {
	return &PlayCmd{
		name:        "play",
		group:       "player",
		summary:     "Play a track, playlist, or search for something on YouTube",
		description: "The motherload of commands!\n\tIf no arguments are provided, then it'll resume any paused track.\n\tIf provided a URL, then it'll attempt to directly play it. Right now, only YouTube and SoundCloud are supported.\n\tIf a non-URL is provided, it'll search YouTube with your query, and provide the first five options. There, you can use `play 1` to select the first option, for example.",
		usage:       "play [url or query]",
		aliases:     []string{"p", "load"},
		voiceOnly:   true,
	}
}

func (cmd *PlayCmd) GetName() string {
	return cmd.name
}

func (cmd *PlayCmd) GetGroup() string {
	return cmd.group
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
		if err != nil && !msg.IsErrorMessage(err, msg.NoLoadedTrack) {
			return err
		} else if status == enum.StatusSuccess {
			deps.Messenger.SendSimpleMessage("Resuming.")
			return nil
		}

		deps.Messenger.SendSimpleMessage("Bruh where's your song??")
		return nil
	}

	// Selecting a search result
	num, err := strconv.Atoi(opts[0])
	if err == nil {
		if num < 1 || num > deps.MusicPlayer.GetSearchResultsLength() {
			deps.Messenger.SendSimpleMessage(fmt.Sprintf("Selected thin air. Try a number between 1 and %d.", deps.MusicPlayer.GetSearchResultsLength()))
			return nil
		}

		client := deps.Event.Client()
		err = deps.MusicPlayer.JoinVoiceChannel(&client, deps.Event.Message.Author.ID)
		if err != nil {
			if msg.IsErrorMessage(err, msg.VoiceStateNotFound) {
				deps.Messenger.SendSimpleMessage("Dude you're not in a voice channel... get in one I can see!")
				return nil
			}

			return err
		}

		track := deps.MusicPlayer.GetSearchResult(num - 1)
		status, err := deps.MusicPlayer.Load(*track)
		if err != nil {
			return err
		}

		if status == enum.StatusQueued {
			deps.Messenger.SendSimpleMessage(fmt.Sprintf("Queued *%s* by **%s**.", track.Title(), track.Author()))
			return nil
		}

		deps.Messenger.SendSimpleMessage(fmt.Sprintf("Now playing *%s* by **%s**.", track.Title(), track.Author()))
		return nil
	}

	// If we have a search query, then we'll need to put the string opts back together into one.
	song := opts[0]
	_, err = url.ParseRequestURI(song)
	if err != nil {
		song = fmt.Sprintf("%s:%s", lavalink.SearchTypeYouTube, strings.Join(opts, " "))
	} else {
		client := deps.Event.Client()
		err = deps.MusicPlayer.JoinVoiceChannel(&client, deps.Event.Message.Author.ID)
		if err != nil {
			if msg.IsErrorMessage(err, msg.VoiceStateNotFound) {
				deps.Messenger.SendSimpleMessage("Dude you're not in a voice channel... get in one I can see!")
				return nil
			}

			return err
		}
	}

	status, track, tracksQueued, err := deps.MusicPlayer.ProcessQuery(song)
	if err != nil {
		if msg.IsErrorMessage(err, msg.NoResultsFound) {
			deps.Messenger.SendSimpleMessage("No results. Try some other keywords? Such as OFFICIAL, FEATURING, ft., THE TRUTH ABOUT, IS A FRAUD, or CHARLIE")
			return nil
		}

		return err
	}

	switch status {
	case enum.StatusQueued:
		deps.Messenger.SendSimpleMessage(fmt.Sprintf("Queued *%s* by **%s**.", track.Title(), track.Author()))
	case enum.StatusQueuedList:
		if tracksQueued == 0 {
			deps.Messenger.SendSimpleMessage("Queued nothing. What the...?")
		} else if tracksQueued == 1 {
			deps.Messenger.SendSimpleMessage("Queued **1** additional track. Just a one hit wonder, huh?")
		} else {
			deps.Messenger.SendSimpleMessage(fmt.Sprintf("Queued **%d** additional tracks.", tracksQueued))
		}
	case enum.StatusPlayingAndQueuedList:
		deps.Messenger.SendSimpleMessage(fmt.Sprintf("Now playing *%s* by **%s**.\nQueued **1** additional track.", track.Title(), track.Author()))
	case enum.StatusSearchSuccess:
		builder := strings.Builder{}
		builder.WriteString("Search Results:\n\n")

		for i, result := range deps.MusicPlayer.GetSearchResults() {
			builder.WriteString(fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", i+1, result.Title(), result.Author(), result.Length().String()))
		}

		builder.WriteString(fmt.Sprintf("\nUse `%splay n` to pick a track to play. Results will be available until the next query.", deps.Prefix))

		deps.Messenger.SendSimpleMessage(builder.String())
	default:
		deps.Messenger.SendSimpleMessage(fmt.Sprintf("Now playing *%s* by **%s**.", track.Title(), track.Author()))
	}

	return nil
}
