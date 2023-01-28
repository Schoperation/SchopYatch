package main

import (
	"context"
	"log"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

type SchopYatch struct {
	Client bot.Client
	Config YatchConfig
}

func NewSchopYatchBot(config YatchConfig) *SchopYatch {
	return &SchopYatch{
		Config: config,
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

	var message string
	if event.Message.Content == "ping" {
		message = "pong"
	} else if event.Message.Content == "pong" {
		message = "ping"
	}

	if message != "" {
		_, err := event.Client().Rest().CreateMessage(event.ChannelID, discord.NewMessageCreateBuilder().SetContent(message).Build())
		if err != nil {
			log.Printf("Error sending the message: %v", err)
		}
	}
}
