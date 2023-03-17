package bot

import (
	"context"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgolink/disgolink"
	"github.com/disgoorg/disgolink/lavalink"
)

func createClient(token string, eventListeners ...bot.EventListener) (bot.Client, error) {
	client, err := disgo.New(token,
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
		bot.WithEventListeners(eventListeners...),
	)

	if err != nil {
		return nil, err
	}

	return client, nil
}

func createLavalinkConn(client bot.Client, config YatchConfig) (disgolink.Link, error) {
	link := disgolink.New(client)
	_, err := link.AddNode(context.TODO(), lavalink.NodeConfig{
		Name:        "schopyatch",
		Host:        config.LavalinkHost,
		Port:        config.LavalinkPort,
		Password:    config.LavalinkPassword,
		Secure:      config.LavalinkSecure,
		ResumingKey: "",
	})

	if err != nil {
		return nil, err
	}

	return link, nil
}
