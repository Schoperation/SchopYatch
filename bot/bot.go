package bot

import (
	"context"
	"fmt"
	"log"
	"schoperation/schopyatch/command"
	"schoperation/schopyatch/music_player"
	"schoperation/schopyatch/util"
	"strings"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgolink/disgolink"
	"github.com/disgoorg/snowflake/v2"
)

type SchopYatch struct {
	Client   bot.Client
	Config   YatchConfig
	Lavalink disgolink.Link
	commands map[string]command.Command
	players  map[snowflake.ID]*music_player.MusicPlayer
	version  string
}

func NewSchopYatchBot(config YatchConfig, version string) (SchopYatch, error) {
	schopYatch := SchopYatch{
		Config:   config,
		commands: command.GetCommandsAndAliasesAsMap(),
		players:  make(map[snowflake.ID]*music_player.MusicPlayer),
		version:  version,
	}

	client, err := createClient(config.Token,
		bot.NewListenerFunc(schopYatch.OnReady),
		bot.NewListenerFunc(schopYatch.OnGuildJoin),
		bot.NewListenerFunc(schopYatch.OnMessageCreate),
	)
	if err != nil {
		return SchopYatch{}, err
	}

	lavalink, err := createLavalinkConn(client, config)
	if err != nil {
		return SchopYatch{}, err
	}

	schopYatch.Client = client
	schopYatch.Lavalink = lavalink

	return schopYatch, nil
}

func (sy *SchopYatch) getPlayerByGuildId(guildId snowflake.ID) *music_player.MusicPlayer {
	player, exists := sy.players[guildId]
	if !exists {
		return nil
	}

	return player
}

func (sy *SchopYatch) OnReady(event *events.Ready) {
	err := event.Client().SetPresence(context.TODO(), gateway.WithListeningActivity("an Ace Attorney OST"))
	if err != nil {
		log.Fatalf("Error setting presence: %v", err)
	}

	log.Printf("SchopYatch v%s is up and running!", sy.version)
}

func (sy *SchopYatch) OnGuildJoin(event *events.GuildJoin) {
	guildId := event.GuildID

	sy.players[guildId] = music_player.NewMusicPlayer(guildId, sy.Lavalink)
}

func (sy *SchopYatch) OnMessageCreate(event *events.MessageCreate) {
	if event.Message.Author.Bot {
		return
	}

	message := event.Message.Content

	if !strings.HasPrefix(message, sy.Config.Prefix) {
		return
	}

	message = strings.Replace(message, sy.Config.Prefix, "", 1)

	splitMessage := strings.Split(message, " ")
	cmd, exists := sy.commands[strings.ToLower(splitMessage[0])]
	if !exists {
		return
	}

	if cmd.IsVoiceOnlyCmd() {
		userVoiceState, exists := sy.Client.Caches().VoiceState(*event.GuildID, event.Message.Author.ID)
		if !exists {
			util.SendSimpleMessage(sy.Client, event.ChannelID, "Dude you're not in a voice channel... get in one I can see!")
			return
		}

		botVoiceState, exists := sy.Client.Caches().VoiceState(*event.GuildID, sy.Client.ID())
		if !exists && cmd.GetName() != "join" && cmd.GetName() != "play" {
			util.SendSimpleMessage(sy.Client, event.ChannelID, fmt.Sprintf("Dude I'm not in a voice channel... use either `%sjoin` or `%splay` to summon me.", sy.Config.Prefix, sy.Config.Prefix))
			return
		}

		if exists {
			if userVoiceState.ChannelID.String() != botVoiceState.ChannelID.String() {
				util.SendSimpleMessage(sy.Client, event.ChannelID, "It would appear that you're in a different channel.")
				return
			}
		}
	}

	player := sy.getPlayerByGuildId(*event.GuildID)
	if player == nil {
		log.Printf("Hol' up, there's no initialized music player for your server?")
		return
	}

	err := cmd.Execute(command.CommandDependencies{
		Client:      &sy.Client,
		Event:       event,
		MusicPlayer: player,
		Lavalink:    &sy.Lavalink,
		Prefix:      sy.Config.Prefix,
	}, splitMessage[1:]...)

	if err != nil {
		log.Printf("Error occurred running the %s command: %v", cmd.GetName(), err)
		return
	}
}
