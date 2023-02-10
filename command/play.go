package command

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"schoperation/schopyatch/util"
	"strconv"
	"strings"

	"github.com/disgoorg/disgolink/lavalink"
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
		if deps.MusicPlayer.Player.Paused() {
			return resume(deps)
		}

		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Bruh where's your song??")
		return nil
	}

	num, err := strconv.Atoi(opts[0])
	if err == nil {
		if num < 1 || num > deps.MusicPlayer.SearchResults().Length() {
			util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, fmt.Sprintf("Selected thin air. Try a number between 1 and %d.", deps.MusicPlayer.SearchResults().Length()))
			return nil
		}

		cmd.playTrack(deps, *deps.MusicPlayer.SearchResults().GetTrack(num - 1))
		deps.MusicPlayer.SearchResults().Clear()
		return nil
	}

	// If we have a query, then we'll need to put the string opts back together into one.
	song := opts[0]
	_, err = url.ParseRequestURI(song)
	if err != nil {
		song = fmt.Sprintf("%s:%s", lavalink.SearchTypeYoutube, strings.Join(opts, " "))
	}

	err = (*deps.Lavalink).BestRestClient().LoadItemHandler(context.TODO(), song, lavalink.NewResultHandler(
		func(track lavalink.AudioTrack) {
			cmd.playTrack(deps, track)
		},
		func(playlist lavalink.AudioPlaylist) {
			cmd.playList(deps, playlist)
		},
		func(tracks []lavalink.AudioTrack) {
			cmd.search(deps, tracks)
		},
		func() {
			util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Track not found. Make sure the URL is correct, or try searching something else...")
		},
		func(ex lavalink.FriendlyException) {
			log.Printf("Lavalink error: %s", ex.Message)
		},
	))

	return err
}

func (cmd *PlayCmd) playTrack(deps CommandDependencies, track lavalink.AudioTrack) {
	err := joinVoiceChannel(deps)
	if err != nil {
		log.Printf("Couldn't join voice channel: %v", err)
		return
	}

	if deps.MusicPlayer.GotDisconnected {
		deps.MusicPlayer.RecreatePlayer(*deps.Lavalink)
		deps.MusicPlayer.GotDisconnected = false
	}

	if deps.MusicPlayer.Player.PlayingTrack() == nil {
		err = deps.MusicPlayer.Player.Play(track)
		if err != nil {
			log.Printf("%v", err)
			util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "For some reason I can't play this... might be some dumb age restriction?")
			return
		}

		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, fmt.Sprintf("Now playing *%s* by **%s**.", track.Info().Title, track.Info().Author))
		return
	}

	deps.MusicPlayer.Queue.Enqueue(track)
	util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, fmt.Sprintf("Added *%s* by **%s** to the queue.", track.Info().Title, track.Info().Author))
}

func (cmd *PlayCmd) playList(deps CommandDependencies, playlist lavalink.AudioPlaylist) {
	if len(playlist.Tracks()) == 0 {
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Empty playlist. Try again...?")
		return
	}

	cmd.playTrack(deps, playlist.Tracks()[0])

	if len(playlist.Tracks()) == 1 {
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Just a one hit wonder, huh?")
		return
	}

	deps.MusicPlayer.Queue.EnqueueList(playlist.Tracks()[1:])
	util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, fmt.Sprintf("Added **%d** additional tracks from playlist **%s** to the queue.", len(playlist.Tracks()[1:]), playlist.Name()))
}

func (cmd *PlayCmd) search(deps CommandDependencies, tracks []lavalink.AudioTrack) {
	if len(tracks) == 0 {
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "No results. Try some other keywords? Such as OFFICIAL, FEATURING, ft., THE TRUTH ABOUT, IS A FRAUD, or CHARLIE")
		return
	}

	builder := strings.Builder{}
	builder.WriteString("Search Results:\n\n")

	rangeLimit := deps.MusicPlayer.SearchResults().MaxLength()
	if rangeLimit > len(tracks) {
		rangeLimit = len(tracks)
	}

	deps.MusicPlayer.SearchResults().Clear()

	for i := 0; i < rangeLimit; i++ {
		builder.WriteString(fmt.Sprintf("`%02d` - *%s* by **%s** `[%s]`\n", i+1, tracks[i].Info().Title, tracks[i].Info().Author, tracks[i].Info().Length))
		deps.MusicPlayer.SearchResults().AddTrack(tracks[i])
	}

	builder.WriteString(fmt.Sprintf("\nUse `%splay n` to pick a track to play.", deps.Prefix))
	util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, builder.String())
}
