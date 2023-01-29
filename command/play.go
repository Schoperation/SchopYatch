package command

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"schoperation/schopyatch/util"
	"strings"

	"github.com/disgoorg/disgolink/lavalink"
)

type PlayCmd struct {
	Name        string
	Summary     string
	Description string
	Usage       string
	Aliases     []string
}

func NewPlayCmd() Command {
	return &PlayCmd{
		Name:        "play",
		Summary:     "Play a track or playlist",
		Description: "Plays a track on the bot",
		Usage:       "play <required> [optional]",
		Aliases:     []string{"p", "resume"},
	}
}

func (cmd *PlayCmd) GetName() string {
	return cmd.Name
}

func (cmd *PlayCmd) GetSummary() string {
	return cmd.Summary
}

func (cmd *PlayCmd) GetDescription() string {
	return cmd.Description
}

func (cmd *PlayCmd) GetUsage() string {
	return cmd.Usage
}

func (cmd *PlayCmd) GetAliases() []string {
	return cmd.Aliases
}

func (cmd *PlayCmd) Execute(deps CommandDependencies, opts ...string) error {

	if len(opts) == 0 {
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Bruh where's your song")
		return nil
	}

	// If we have a query, then we'll need to put the string opts back together into one.
	song := opts[0]
	_, err := url.ParseRequestURI(song)
	if err != nil {
		song = fmt.Sprintf("%s:%s", lavalink.SearchTypeYoutube, strings.Join(opts, " "))
	}

	err = (*deps.Lavalink).BestRestClient().LoadItemHandler(context.TODO(), song, lavalink.NewResultHandler(
		func(track lavalink.AudioTrack) {
			// Loaded a single track
			cmd.playTrack(deps, track)
		},
		func(playlist lavalink.AudioPlaylist) {
			// Loaded a playlist
			util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Playlist loading not implemented yet, sory :((")
		},
		func(tracks []lavalink.AudioTrack) {
			// Loaded a search result
			util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Searching not implemented yet, sory :((")
		},
		func() {
			// nothing matching the query found
			util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Track not found. Make sure the URL is correct, or try searching something else...")
		},
		func(ex lavalink.FriendlyException) {
			// something went wrong while loading the track
			log.Printf("Lavalink error: %s", ex.Message)
		},
	))

	return err
}

func (cmd *PlayCmd) playTrack(deps CommandDependencies, track lavalink.AudioTrack) {

	voiceState, exists := (*deps.Client).Caches().VoiceState(*deps.Event.GuildID, deps.Event.Message.Author.ID)
	if !exists {
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Dude you're not in a voice channel... get in one I can see!")
		return
	}

	err := (*deps.Client).UpdateVoiceState(context.TODO(), *deps.Event.GuildID, voiceState.ChannelID, false, true)
	if err != nil {
		util.SendSimpleMessage(*deps.Client, deps.Event.ChannelID, "Cannot connect to your channel... do I have permission?")
		log.Printf("%v", err)
		return
	}

	player := (*deps.Lavalink).Player(*deps.Event.GuildID)

	err = player.SetVolume(42)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	err = player.Play(track)
	if err != nil {
		log.Printf("%v", err)
		return
	}

}
