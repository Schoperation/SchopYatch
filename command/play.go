package command

import (
	"context"
	"log"

	"github.com/disgoorg/disgolink/lavalink"
	"github.com/disgoorg/snowflake/v2"
)

type PlayCmd struct {
	Name        string
	Description string
}

func NewPlayCmd() Command {
	return &PlayCmd{
		Name:        "play",
		Description: "Plays a track on the bot",
	}
}

func (cmd *PlayCmd) GetName() string {
	return cmd.Name
}

func (cmd *PlayCmd) GetDescription() string {
	return cmd.Description
}

func (cmd *PlayCmd) Execute(deps CommandDependencies, opts ...string) error {
	//query := "ytsearch:Rick Astley - Never Gonna Give You Up"
	query := "https://www.youtube.com/watch?v=ez-T4UD-bzs"
	//query := "https://soundcloud.com/milulumilu/phoenix-wright-ace-attorney-11"

	err := (*deps.Lavalink).BestRestClient().LoadItemHandler(context.TODO(), query, lavalink.NewResultHandler(
		func(track lavalink.AudioTrack) {
			// Loaded a single track
			//track2 = track
			cmd.play(deps, track)
		},
		func(playlist lavalink.AudioPlaylist) {
			// Loaded a playlist
		},
		func(tracks []lavalink.AudioTrack) {
			// Loaded a search result
		},
		func() {
			// nothing matching the query found
		},
		func(ex lavalink.FriendlyException) {
			// something went wrong while loading the track
			log.Printf("Lavalink error: %s", ex.Message)
		},
	))

	return err
}

func (cmd *PlayCmd) play(deps CommandDependencies, track lavalink.AudioTrack) {

	channelId, err := snowflake.Parse("12345")
	if err != nil {
		log.Printf("%v", err)
	}
	err = (*deps.Client).UpdateVoiceState(context.TODO(), *deps.Event.GuildID, &channelId, false, true)
	if err != nil {
		log.Printf("%v", err)
	}

	log.Printf("%s", track.Info().Title)

	player := (*deps.Lavalink).Player(*deps.Event.GuildID)
	err = player.Play(track)
	if err != nil {
		log.Printf("%v", err)
	}
}
