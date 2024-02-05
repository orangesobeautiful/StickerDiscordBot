package usecase

import (
	"context"

	"backend/app/domain"
	"backend/app/ent/schema"
	"backend/app/pkg/hserr"

	"github.com/sashabaranov/go-openai"
	"golang.org/x/xerrors"
)

// GPT3Dot5Turbo0125 OpenAI GPT-3.5-turbo-0125 model
//
//	openai package 尚未更新新模型，這裡用自定義處理
const GPT3Dot5Turbo0125 = "gpt-3.5-turbo-0125"

var _ domain.ChatUsecase = (*chatUsecase)(nil)

type chatUsecase struct {
	chatRepo domain.ChatRepository

	discordGuildUsecase domain.DiscordGuildUsecase

	openaiClient *openai.Client
}

func New(
	chatRepo domain.ChatRepository,
	discordGuildUsecase domain.DiscordGuildUsecase,
	openaiClient *openai.Client,
) domain.ChatUsecase {
	return &chatUsecase{
		chatRepo: chatRepo,

		discordGuildUsecase: discordGuildUsecase,

		openaiClient: openaiClient,
	}
}

func (u *chatUsecase) Chat(
	ctx context.Context, guildID string, message string,
) (replyMessage string, err error) {
	chatroom, err := u.discordGuildUsecase.GetGuildActivateChatroom(ctx, guildID)
	if err != nil {
		return "", xerrors.Errorf("get guild activate chatroom: %w", err)
	}

	const maxChatHistoryRefrenceAmount = 5
	refChatHistoriesResult, err := u.chatRepo.ListChatHistory(ctx, chatroom.ID, 0, maxChatHistoryRefrenceAmount)
	if err != nil {
		return "", xerrors.Errorf("list chat histories: %w", err)
	}
	chatMessages := createChatMessagesByChatHistories(message, refChatHistoriesResult)

	const chatModel = GPT3Dot5Turbo0125
	chatResuest := openai.ChatCompletionRequest{
		Model:    chatModel,
		Messages: chatMessages,
	}
	chatResp, err := u.openaiClient.CreateChatCompletion(ctx, chatResuest)
	if err != nil {
		return "", hserr.NewInternalError(err, "create chat completion")
	}
	replyMessage = chatResp.Choices[0].Message.Content
	promptPrice, completePrice, err := cacluteOpenaiChatUsagePrice(chatModel, chatResp.Usage)
	if err != nil {
		return "", xerrors.Errorf("calculate openai chat usage price: %w", err)
	}

	_, err = u.chatRepo.CreateOpenaiChatHistory(ctx,
		chatroom.ID,
		chatResp.Model,
		message,
		replyMessage,
		openaiChatMessagesToEntChatMessages(chatMessages),
		&chatResuest,
		&chatResp,
		promptPrice,
		completePrice,
	)
	if err != nil {
		return "", xerrors.Errorf("create chat history: %w", err)
	}

	return replyMessage, nil
}

func createChatMessagesByChatHistories(
	message string, refReverseChatHistoriesResult domain.ListChatHistoriesResult,
) (messages []openai.ChatCompletionMessage) {
	refHistoryMessages := reverseChatHistoriesResultToOpenaiChatCompletionMessages(refReverseChatHistoriesResult)

	chatMessages := refHistoryMessages
	chatMessages = append(chatMessages,
		openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: message,
		},
	)

	return chatMessages
}

func reverseChatHistoriesResultToOpenaiChatCompletionMessages(
	reverseChatHistoriesResult domain.ListChatHistoriesResult,
) (messages []openai.ChatCompletionMessage) {
	items := reverseChatHistoriesResult.GetItems()
	for i := len(items) - 1; i >= 0; i-- {
		chatHistory := items[i]
		messages = append(messages,
			openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: chatHistory.RequestMessage,
			},
			openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: chatHistory.ReplyMessage,
			},
		)
	}

	return messages
}

func openaiChatMessagesToEntChatMessages(
	openaiChatMessages []openai.ChatCompletionMessage,
) (entChatMessages []schema.ChatMessage) {
	entChatMessages = make([]schema.ChatMessage, len(openaiChatMessages))

	for i, openaiChatMessage := range openaiChatMessages {
		entChatMessages[i] = schema.ChatMessage{
			Type:    openaiChatMessage.Role,
			Content: openaiChatMessage.Content,
		}
	}

	return entChatMessages
}
