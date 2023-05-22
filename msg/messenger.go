package msg

import (
	"log"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

type Messenger struct {
	client  *bot.Client
	channel snowflake.ID
}

func NewMessenger(client *bot.Client) Messenger {
	return Messenger{
		client: client,
	}
}

func (msgr *Messenger) SetChannel(channel snowflake.ID) {
	msgr.channel = channel
}

func (msgr *Messenger) SendSimpleMessage(msg string) {
	_, err := (*msgr.client).Rest().CreateMessage(msgr.channel, discord.NewMessageCreateBuilder().SetContent(msg).ClearEmbeds().Build())
	if err != nil {
		log.Printf("Couldn't send a message: %v", err)
	}
}
