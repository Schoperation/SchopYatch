package main

import (
	"context"
	"log"
	"schoperation/schopyatch/command"
	"strings"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

type SchopYatch struct {
	Client   bot.Client
	Config   YatchConfig
	Commands map[string]command.Command
}

func NewSchopYatchBot(config YatchConfig, commands []command.Command) *SchopYatch {
	var mappedCommands = make(map[string]command.Command)

	for _, command := range commands {
		mappedCommands[command.GetName()] = command
	}

	return &SchopYatch{
		Config:   config,
		Commands: mappedCommands,
	}
}

func (sy *SchopYatch) SetupClient() error {
	var err error
	sy.Client, err = disgo.New(sy.Config.Token,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(
				gateway.IntentGuildMessages,
				gateway.IntentMessageContent,
			),
		),
		bot.WithEventListenerFunc(sy.OnReady),
		bot.WithEventListenerFunc(sy.OnMessageCreate),
	)
	if err != nil {
		return err
	}

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
	cmd, exists := sy.Commands[splitMessage[0]]
	if !exists {
		return
	}

	cmd.Execute(event, splitMessage[1:]...)
}
