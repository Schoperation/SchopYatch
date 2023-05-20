package command_tests

import (
	"context"

	"github.com/disgoorg/log"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/httpserver"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/disgo/sharding"
	"github.com/disgoorg/disgo/voice"
	"github.com/disgoorg/snowflake/v2"
)

type fakeBotClient struct{}

func NewFakeBotClient() bot.Client {
	return fakeBotClient{}
}

func NewFakeMessageCreateEvent() *events.MessageCreate {
	return &events.MessageCreate{
		GenericMessage: &events.GenericMessage{
			GenericEvent: events.NewGenericEvent(NewFakeBotClient(), 1, 1),
		},
	}
}

func (client fakeBotClient) Logger() log.Logger {
	return log.New(1)
}

func (client fakeBotClient) Close(ctx context.Context) {

}

func (client fakeBotClient) Token() string {
	return ""
}

func (client fakeBotClient) ApplicationID() snowflake.ID {
	return 1
}

func (client fakeBotClient) ID() snowflake.ID {
	return 1
}

func (client fakeBotClient) Caches() cache.Caches {
	return nil
}

func (client fakeBotClient) Rest() rest.Rest {
	return nil
}

func (client fakeBotClient) AddEventListeners(listeners ...bot.EventListener) {

}

func (client fakeBotClient) RemoveEventListeners(listeners ...bot.EventListener) {

}

func (client fakeBotClient) EventManager() bot.EventManager {
	return nil
}

func (client fakeBotClient) VoiceManager() voice.Manager {
	return nil
}

func (client fakeBotClient) OpenGateway(ctx context.Context) error {
	return nil
}

func (client fakeBotClient) Gateway() gateway.Gateway {
	return nil
}

func (client fakeBotClient) HasGateway() bool {
	return false
}

func (client fakeBotClient) OpenShardManager(ctx context.Context) error {
	return nil
}

func (client fakeBotClient) ShardManager() sharding.ShardManager {
	return nil
}

func (client fakeBotClient) HasShardManager() bool {
	return false
}

func (client fakeBotClient) Shard(guildID snowflake.ID) (gateway.Gateway, error) {
	return nil, nil
}

func (client fakeBotClient) UpdateVoiceState(ctx context.Context, guildID snowflake.ID, channelID *snowflake.ID, selfMute bool, selfDeaf bool) error {
	return nil
}

func (client fakeBotClient) RequestMembers(ctx context.Context, guildID snowflake.ID, presence bool, nonce string, userIDs ...snowflake.ID) error {
	return nil
}

func (client fakeBotClient) RequestMembersWithQuery(ctx context.Context, guildID snowflake.ID, presence bool, nonce string, query string, limit int) error {
	return nil
}

func (client fakeBotClient) SetPresence(ctx context.Context, opts ...gateway.PresenceOpt) error {
	return nil
}

func (client fakeBotClient) SetPresenceForShard(ctx context.Context, shardId int, opts ...gateway.PresenceOpt) error {
	return nil
}

func (client fakeBotClient) MemberChunkingManager() bot.MemberChunkingManager {
	return nil
}

func (client fakeBotClient) OpenHTTPServer() error {
	return nil
}

func (client fakeBotClient) HTTPServer() httpserver.Server {
	return nil
}

func (client fakeBotClient) HasHTTPServer() bool {
	return false
}
