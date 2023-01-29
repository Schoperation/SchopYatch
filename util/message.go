package util

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

func SendSimpleMessage(client bot.Client, channel snowflake.ID, message string) error {
	_, err := client.Rest().CreateMessage(channel, discord.NewMessageCreateBuilder().SetContent(message).Build())
	if err != nil {
		return err
	}

	return nil
}
