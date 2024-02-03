package repository

import (
	"context"

	"backend/app/domain"
	"backend/app/ent"
	"backend/app/ent/chathistory"
	"backend/app/ent/chathistorydetail"
	"backend/app/ent/chatroom"
	"backend/app/ent/schema"
	"backend/app/pkg/hserr"

	"github.com/shopspring/decimal"
)

var _ domain.ChatRepository = (*chatRepository)(nil)

type chatRepository struct {
	*domain.BaseEntRepo
}

func New(client *ent.Client) domain.ChatRepository {
	bRepo := domain.NewBaseEntRepo(client)

	return &chatRepository{
		BaseEntRepo: bRepo,
	}
}

func (r *chatRepository) CreateChatHistory(ctx context.Context,
	chatroomID int,
	model, requestMessage, replyMessage string,
	fullRequestMessage []schema.ChatMessage,
	args schema.ChatMessageRequestArgument,
	promptTokens, completionTokens uint,
	promptPrice, completionPrice decimal.Decimal,
) (chatHistoryID int, err error) {
	newChatHistory, err := r.GetEntClient(ctx).ChatHistory.
		Create().
		SetChatroomID(chatroomID).
		SetRequestMessage(requestMessage).
		SetReplyMessage(replyMessage).
		Save(ctx)
	if err != nil {
		return 0, hserr.NewInternalError(err, "create chat history")
	}

	_, err = r.GetEntClient(ctx).ChatHistoryDetail.
		Create().
		SetRefID(newChatHistory.ID).
		SetModel(model).
		SetFullRequestMessage(fullRequestMessage).
		SetRequestArgument(args).
		SetPromptTokens(promptTokens).
		SetCompletionTokens(completionTokens).
		SetPromptPrice(promptPrice).
		SetCompletionPrice(completionPrice).
		Save(ctx)
	if err != nil {
		return 0, hserr.NewInternalError(err, "create chat history detail")
	}

	return newChatHistory.ID, nil
}

func (r *chatRepository) ListChatHistory(
	ctx context.Context, chatroomID, offset, limit int,
) (result domain.ListChatHistoriesResult, err error) {
	queryFilter := r.GetEntClient(ctx).ChatHistory.
		Query().
		Where(
			chathistory.HasChatroomWith(chatroom.ID(chatroomID)),
		).
		Order(ent.Desc(chatroom.FieldCreatedAt)).
		Offset(offset).
		Limit(limit)

	total, err := queryFilter.Clone().Count(ctx)
	if err != nil {
		return result, hserr.NewInternalError(err, "query chat history count")
	}

	chatHistories, err := queryFilter.All(ctx)
	if err != nil {
		return domain.ListChatHistoriesResult{}, hserr.NewInternalError(err, "list chat history")
	}

	result = domain.NewListResult(total, chatHistories)
	return result, nil
}

func (r *chatRepository) FindChatHistoryDetailByChatHistoryID(
	ctx context.Context, chatHistoryID int,
) (detail *ent.ChatHistoryDetail, err error) {
	detail, err = r.GetEntClient(ctx).ChatHistoryDetail.
		Query().
		Where(
			chathistorydetail.HasRefWith(chathistory.ID(chatHistoryID)),
		).
		Only(ctx)
	if err != nil {
		return nil, hserr.NewInternalError(err, "find chat history detail")
	}

	return detail, nil
}
