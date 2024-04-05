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

	CreateRAGReferencePool(ctx context.Context, guildID, name, description string) (id int, err error)
	ListRAGReferencePools(ctx context.Context, guildID string, limit, offset int) (result ListRAGReferencePoolsResult, err error)
	GetRAGReferencePoolWithGuildByID(ctx context.Context, ragReferencePoolID int) (ragReferencePool *ent.RAGReferencePool, err error)
	SearchRAGReferencePoolText(
		ctx context.Context, ragReferencePoolID []int, reqEmbedding []float32, topK uint) (result []string, err error)

	CreateRAGReferenceText(
		ctx context.Context, ragReferencePoolID int, text string, embedContent []float32) (id int, err error)
	GetRAGReferenceTextContent(ctx context.Context, ragReferenceTextID int) (content string, err error)
	ListRAGReferenceTexts(ctx context.Context, ragReferencePoolID int, limit, offset int) (result ListRAGReferenceTextsResult, err error)
	GetRAGReferenceTextWithGuildByID(ctx context.Context, ragReferenceTextID int) (ragReferenceText *ent.RAGReferenceText, err error)
	DeleteRAGReferenceText(ctx context.Context, ragReferenceTextID int) (err error)
}

type ChatUsecase interface {
	Chat(ctx context.Context, guildID string, message string) (replyMessage string, err error)

	CreateRAGReferenceText(ctx context.Context, ragReferencePoolID int, text string) (id int, err error)
	ListRAGReferenceTexts(ctx context.Context, ragReferencePoolID int, limit, offset int) (result ListRAGReferenceTextsResult, err error)
	DeleteRAGReferenceText(ctx context.Context, ragReferenceTextID int) (err error)
}
