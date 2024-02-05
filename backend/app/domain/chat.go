package domain

import (
	"context"

	"backend/app/ent"
	"backend/app/ent/schema"

	"github.com/sashabaranov/go-openai"
	"github.com/shopspring/decimal"
)

type ListChatHistoriesResult = ListResult[*ent.ChatHistory]

type ChatRepository interface {
	BaseEntRepoInterface
	CreateOpenaiChatHistory(ctx context.Context,
		chatroomID int,
		model, requestMessage, replyMessage string,
		fullRequestMessage []schema.ChatMessage,
		req *openai.ChatCompletionRequest, resp *openai.ChatCompletionResponse,
		promptPrice, completionPrice decimal.Decimal,
	) (chatHistoryID int, err error)
	ListChatHistory(ctx context.Context, chatroomID, offset, limit int) (result ListChatHistoriesResult, err error)
	FindChatHistoryDetailByChatHistoryID(ctx context.Context, chatHistoryID int) (detail *ent.ChatHistoryDetail, err error)
}

type ChatUsecase interface {
	Chat(ctx context.Context, guildID string, message string) (replyMessage string, err error)
}
