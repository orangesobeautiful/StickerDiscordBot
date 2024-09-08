package delivery

import (
	"time"

	"backend/app/domain"
	"backend/app/ent"
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

type ginCreateRAGReferencePoolTextsReq struct {
	Text string `json:"text" binding:"required"`
}

type createRAGReferencePoolTextsResp struct {
	ID uint64 `json:"id"`
}

type ginListRAGReferencePoolTextsReq struct {
	Page int `form:"page" binding:"required,gte=1"`

	Limit int `form:"limit" binding:"required,gte=1,lte=30"`
}

type listRAGReferencePoolTextsResp struct {
	TotalCount int `json:"total_count"`

	Items []*ragReferencePoolText `json:"texts"`
}

func newListRAGReferencePoolTextsResp(
	listResult domain.ListRAGReferenceTextsResult,
) *listRAGReferencePoolTextsResp {
	return &listRAGReferencePoolTextsResp{
		TotalCount: listResult.GetTotal(),

		Items: newRAGReferencePoolTextsFromEnt(listResult.GetItems()),
	}
}

type ragReferencePoolText struct {
	ID uint64 `json:"id"`

	Text string `json:"text"`

	CreatedAt time.Time `json:"created_at"`
}

func newRAGReferencePoolTextsFromEnt(e []*ent.RAGReferenceText) []*ragReferencePoolText {
	result := make([]*ragReferencePoolText, 0, len(e))

	for _, v := range e {
		result = append(result, newRAGReferencePoolTextFromEnt(v))
	}

	return result
}

func newRAGReferencePoolTextFromEnt(e *ent.RAGReferenceText) *ragReferencePoolText {
	return &ragReferencePoolText{
		ID: e.ID,

		Text: e.Text,

		CreatedAt: e.CreatedAt,
	}
}
