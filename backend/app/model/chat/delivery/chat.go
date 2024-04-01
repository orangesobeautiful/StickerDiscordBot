package delivery

import (
	"context"

	"backend/app/domain"
	domainresponse "backend/app/domain-response"
	ginauth "backend/app/model/discorduser/gin-auth"
	discordcommand "backend/app/pkg/discord-command"
	"backend/app/pkg/ginext"

	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"
)

type chatController struct {
	auth ginauth.AuthInterface

	rd *domainresponse.DomainResponse

	chatUsecase domain.ChatUsecase

	dcGuildUsecase domain.DiscordGuildUsecase
}

func Initialze(
	apiGroup *gin.RouterGroup, dcCmdRegister discordcommand.Register,
	auth ginauth.AuthInterface,
	rd *domainresponse.DomainResponse,
	chatUsecase domain.ChatUsecase, dcGuildUsecase domain.DiscordGuildUsecase,
) {
	ctrl := chatController{
		auth: auth,

		rd: rd,

		chatUsecase: chatUsecase,

		dcGuildUsecase: dcGuildUsecase,
	}

	ctrl.registerGinRouter(apiGroup)
	ctrl.registerDiscordCommand(dcCmdRegister)
}

func (c *chatController) ginChat(ctx *gin.Context, req *GinChatReq) (*ChatResp, error) {
	user := c.auth.MustGetUserFromContext(ctx)

	chatReq := &ChatReq{
		GuildID: user.GuildID,
		Message: req.Message,
	}
	resp, err := c.chat(ctx, chatReq)
	if err != nil {
		return nil, xerrors.Errorf("chat: %w", err)
	}

	return resp, nil
}

func (c *chatController) discordChat(req *DiscordChatReq) (*ChatResp, error) {
	chatReq := &ChatReq{
		GuildID: req.GuildID,
		Message: req.Message,
	}
	resp, err := c.chat(context.Background(), chatReq)
	if err != nil {
		return nil, xerrors.Errorf("chat: %w", err)
	}

	return resp, nil
}

func (c *chatController) chat(ctx context.Context, req *ChatReq) (*ChatResp, error) {
	replyMessage, err := c.chatUsecase.Chat(ctx, req.GuildID, req.Message)
	if err != nil {
		return nil, xerrors.Errorf("chat: %w", err)
	}

	return &ChatResp{
		Message: replyMessage,
	}, nil
}

func (c *chatController) ginCreateRAGReferencePoolTexts(
	ctx *gin.Context, req *ginCreateRAGReferencePoolTextsReq,
) (*createRAGReferencePoolTextsResp, error) {
	ragReferencePoolID, err := getRAGReferencePoolIDFromContext(ctx)
	if err != nil {
		return nil, xerrors.Errorf("get rag reference pool id from context: %w", err)
	}

	resp, err := c.createRAGReferencePoolTexts(ctx, ragReferencePoolID, req.Text)
	if err != nil {
		return nil, xerrors.Errorf("create rag reference pool texts: %w", err)
	}

	return resp, nil
}

func (c *chatController) createRAGReferencePoolTexts(
	ctx context.Context, ragReferencePoolID int, text string,
) (*createRAGReferencePoolTextsResp, error) {
	id, err := c.chatUsecase.CreateRAGReferenceText(ctx, ragReferencePoolID, text)
	if err != nil {
		return nil, xerrors.Errorf("create rag reference pool texts: %w", err)
	}

	return &createRAGReferencePoolTextsResp{
		ID: id,
	}, nil
}

func (c *chatController) ginListRAGReferencePoolTexts(
	ctx *gin.Context, req ginListRAGReferencePoolTextsReq,
) (*listRAGReferencePoolTextsResp, error) {
	ragReferencePoolID, err := getRAGReferencePoolIDFromContext(ctx)
	if err != nil {
		return nil, xerrors.Errorf("get rag reference pool id from context: %w", err)
	}

	resp, err := c.listRAGReferencePoolTexts(ctx, ragReferencePoolID, req.Page, req.Limit)
	if err != nil {
		return nil, xerrors.Errorf("list rag reference pool texts: %w", err)
	}

	return resp, nil
}

func (c *chatController) listRAGReferencePoolTexts(
	ctx context.Context, ragReferencePoolID, page, limit int,
) (*listRAGReferencePoolTextsResp, error) {
	offset := (page - 1) * limit
	listResult, err := c.chatUsecase.ListRAGReferenceTexts(ctx, ragReferencePoolID, limit, offset)
	if err != nil {
		return nil, xerrors.Errorf("list rag reference pool texts: %w", err)
	}

	return newListRAGReferencePoolTextsResp(listResult), nil
}

func (c *chatController) ginDeleteRAGReferenceText(ctx *gin.Context) (*ginext.EmptyResp, error) {
	ragReferenceTextID, err := getRAGReferenceTextIDFromContext(ctx)
	if err != nil {
		return nil, xerrors.Errorf("get rag reference text id from context: %w", err)
	}

	err = c.deleteRAGReferenceText(ctx, ragReferenceTextID)
	if err != nil {
		return nil, xerrors.Errorf("delete rag reference text: %w", err)
	}

	return nil, nil
}

func (c *chatController) deleteRAGReferenceText(ctx context.Context, ragReferenceTextID int) error {
	err := c.chatUsecase.DeleteRAGReferenceText(ctx, ragReferenceTextID)
	if err != nil {
		return xerrors.Errorf("delete rag reference text: %w", err)
	}

	return nil
}
