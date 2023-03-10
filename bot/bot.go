package bot

import (
	"context"
	"fmt"
	"log"
	"schoperation/schopyatch/command"
	"schoperation/schopyatch/musicplayer"
	"schoperation/schopyatch/util"
	"strings"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgolink/disgolink"
	"github.com/disgoorg/disgolink/lavalink"
	"github.com/disgoorg/snowflake/v2"
)

type SchopYatch struct {
	Client   bot.Client
	Config   YatchConfig
	Commands map[string]command.Command
	guilds   map[snowflake.ID]string
	players  map[snowflake.ID]*musicplayer.MusicPlayer
	Lavalink disgolink.Link
}

func NewSchopYatchBot(config YatchConfig) *SchopYatch {
	return &SchopYatch{
		Config:   config,
		Commands: command.GetCommandsAndAliasesAsMap(),
		guilds:   make(map[snowflake.ID]string),
		players:  make(map[snowflake.ID]*musicplayer.MusicPlayer),
	}
}

func (sy *SchopYatch) SetupClient() error {
	var err error
	sy.Client, err = disgo.New(sy.Config.Token,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuilds,
				gateway.IntentGuildMessages,
				gateway.IntentGuildVoiceStates,
				gateway.IntentMessageContent,
			),
		),
		bot.WithCacheConfigOpts(
			cache.WithCaches(cache.FlagVoiceStates, cache.FlagMembers, cache.FlagChannels, cache.FlagGuilds, cache.FlagRoles),
		),
		bot.WithEventListenerFunc(sy.OnReady),
		bot.WithEventListenerFunc(sy.OnMessageCreate),
		bot.WithEventListenerFunc(sy.OnGuildJoin),
	)

	return err
}

func (sy *SchopYatch) SetupLavalink() error {
	link := disgolink.New(sy.Client)
	_, err := link.AddNode(context.TODO(), lavalink.NodeConfig{
		Name:        "schopyatch",
		Host:        sy.Config.LavalinkHost,
		Port:        sy.Config.LavalinkPort,
		Password:    sy.Config.LavalinkPassword,
		Secure:      sy.Config.LavalinkSecure,
		ResumingKey: "",
	})
	if err != nil {
		return err
	}

	sy.Lavalink = link
	return nil
}

func (sy *SchopYatch) getPlayerByGuildId(guildId snowflake.ID) *musicplayer.MusicPlayer {
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

	log.Printf("SchopYatch is up and running!")
}

func (sy *SchopYatch) OnGuildJoin(event *events.GuildJoin) {
	guildId := event.GuildID

	sy.guilds[guildId] = guildId.String()
	sy.players[guildId] = musicplayer.NewMusicPlayer(guildId, sy.Lavalink)
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
	cmd, exists := sy.Commands[strings.ToLower(splitMessage[0])]
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
