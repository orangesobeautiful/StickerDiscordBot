package delivery

import (
	discordcommand "backend/app/pkg/discord-command"

	"github.com/bwmarrin/discordgo"
)

type GinChatReq struct {
	Message string `json:"message" binding:"required"`
}

type DiscordChatReq struct {
	discordcommand.BaseAuthInteractionCreate `dccmd:"ignore"`

	Message string `dccmd:"name=msg" binding:"required"`
}

type ChatReq struct {
	GuildID string
	Message string
}

var _ discordcommand.DiscordWebhookParamsMarshaler = (*ChatResp)(nil)

type ChatResp struct {
	Message string `json:"message"`
}

func (r *ChatResp) MarshalDiscordWebhookParams() *discordgo.WebhookParams {
	result := new(discordgo.WebhookParams)

	result.Content = r.Message

	return result
}
