package util

import (
	"log"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

func SendSimpleMessage(client bot.Client, channel snowflake.ID, message string) {
	_, err := client.Rest().CreateMessage(channel, discord.NewMessageCreateBuilder().SetContent(message).Build())
	if err != nil {
		log.Printf("Couldn't send a message: %v", err)
	}
}
