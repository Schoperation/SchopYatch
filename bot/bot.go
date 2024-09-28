package bot

import (
	"context"
	"fmt"
	"log"
	"schoperation/schopyatch/command"
	"schoperation/schopyatch/msg"
	"schoperation/schopyatch/music_player"
	"strings"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/snowflake/v2"
)

const SchopYatchVersion = "1.2.1"

type SchopYatch struct {
	Client         bot.Client
	Config         YatchConfig
	LavalinkClient disgolink.Client
	messenger      msg.Messenger
	commands       map[string]command.Command
	players        map[snowflake.ID]*music_player.MusicPlayer
	version        string
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
		bot.NewListenerFunc(schopYatch.OnGuildLeave),
		bot.NewListenerFunc(schopYatch.OnVoiceStateUpdate),
		bot.NewListenerFunc(schopYatch.OnVoiceServerUpdate),
		bot.NewListenerFunc(schopYatch.OnMessageCreate),
	)
	if err != nil {
		return SchopYatch{}, err
	}

	lavalink, err := createLavalinkConn(client, config, &schopYatch.players)
	if err != nil {
		return SchopYatch{}, err
	}

	schopYatch.Client = client
	schopYatch.LavalinkClient = lavalink
	schopYatch.messenger = msg.NewMessenger(&client)

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
	err := event.Client().SetPresence(context.TODO(), gateway.WithListeningActivity("an Ace Attorney OST"))
	if err != nil {
		log.Fatalf("Error setting presence: %v", err)
	}

	guildId := event.GuildID

	sy.players[guildId] = music_player.NewMusicPlayer(guildId, &sy.LavalinkClient)
}

func (sy *SchopYatch) OnGuildLeave(event *events.GuildLeave) {
	delete(sy.players, event.GuildID)
}

func (sy *SchopYatch) OnVoiceStateUpdate(event *events.GuildVoiceStateUpdate) {
	if event.VoiceState.UserID != sy.Client.ApplicationID() {
		return
	}

	sy.LavalinkClient.OnVoiceStateUpdate(context.TODO(), event.VoiceState.GuildID, event.VoiceState.ChannelID, event.VoiceState.SessionID)
}

func (sy *SchopYatch) OnVoiceServerUpdate(event *events.VoiceServerUpdate) {
	sy.LavalinkClient.OnVoiceServerUpdate(context.TODO(), event.GuildID, event.Token, *event.Endpoint)
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

	sy.messenger.SetChannel(event.ChannelID)

	if cmd.IsVoiceOnlyCmd() {
		userVoiceState, exists := sy.Client.Caches().VoiceState(*event.GuildID, event.Message.Author.ID)
		if !exists {
			sy.messenger.SendSimpleMessage("Dude you're not in a voice channel... get in one I can see!")
			return
		}

		botVoiceState, exists := sy.Client.Caches().VoiceState(*event.GuildID, sy.Client.ID())
		if !exists && cmd.GetName() != "join" && cmd.GetName() != "play" {
			sy.messenger.SendSimpleMessage(fmt.Sprintf("Dude I'm not in a voice channel... use either `%sjoin` or `%splay` to summon me.", sy.Config.Prefix, sy.Config.Prefix))
			return
		}

		if exists {
			if userVoiceState.ChannelID.String() != botVoiceState.ChannelID.String() {
				sy.messenger.SendSimpleMessage("It would appear that you're in a different channel.")
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
		Event:       event,
		MusicPlayer: player,
		Messenger:   &sy.messenger,
		Prefix:      sy.Config.Prefix,
		Version:     SchopYatchVersion,
	}, splitMessage[1:]...)

	if err != nil {
		log.Printf("Error occurred running the %s command: %v", cmd.GetName(), err)
		sy.messenger.SendSimpleMessage("Unexpected error occurred. Please try again. If this persists then you might wanna file a bug report...")
		return
	}
}
