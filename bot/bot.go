package bot

import (
	"context"
	"log"
	"schoperation/schopyatch/command"
	"strings"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgolink/disgolink"
	"github.com/disgoorg/disgolink/lavalink"
)

type SchopYatch struct {
	Client   bot.Client
	Config   YatchConfig
	Lavalink disgolink.Link
	Commands map[string]command.Command
}

func NewSchopYatchBot(config YatchConfig) *SchopYatch {
	return &SchopYatch{
		Config:   config,
		Commands: command.GetCommandsAndAliasesAsMap(),
	}
}

func (sy *SchopYatch) SetupClient() error {
	var err error
	sy.Client, err = disgo.New(sy.Config.Token,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuildMessages,
				gateway.IntentMessageContent,
				gateway.IntentGuildVoiceStates,
			),
		),
		bot.WithEventListenerFunc(sy.OnReady),
		bot.WithEventListenerFunc(sy.OnMessageCreate),
	)
	if err != nil {
		return err
	}

	link := disgolink.New(sy.Client)
	link.AddNode(context.TODO(), lavalink.NodeConfig{
		Name:        "schopyatch",
		Host:        "localhost",
		Port:        "2333",
		Password:    sy.Config.LavalinkPassword,
		Secure:      false,
		ResumingKey: "",
	})
	sy.Lavalink = link

	return nil
}

func (sy *SchopYatch) OnReady(event *events.Ready) {
	err := event.Client().SetPresence(context.TODO(), gateway.WithListeningActivity("an Ace Attorney OST"))
	if err != nil {
		log.Fatalf("Error setting presence: %v", err)
	}

	log.Printf("SchopYatch is up and running!")
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

	err := cmd.Execute(command.CommandDependencies{
		Client:   &sy.Client,
		Lavalink: &sy.Lavalink,
		Event:    event,
		Prefix:   sy.Config.Prefix,
	}, splitMessage[1:]...)

	if err != nil {
		log.Printf("Error occurred running the %s command: %v", cmd.GetName(), err)
	}
}
