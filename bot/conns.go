package bot

import (
	"context"
	"fmt"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgolink/v2/disgolink"
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

func createLavalinkConn(client bot.Client, config YatchConfig) (disgolink.Client, error) {
	link := disgolink.New(client.ID())
	_, err := link.AddNode(context.TODO(), disgolink.NodeConfig{
		Name:     "schopyatch",
		Address:  fmt.Sprintf("%s:%s", config.LavalinkHost, config.LavalinkPort),
		Password: config.LavalinkPassword,
		Secure:   config.LavalinkSecure,
	})

	if err != nil {
		return nil, err
	}

	return link, nil
}
