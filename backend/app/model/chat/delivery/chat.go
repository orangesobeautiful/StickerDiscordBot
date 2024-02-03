package delivery

import (
	"context"

	"backend/app/domain"
	domainresponse "backend/app/domain-response"
	ginauth "backend/app/model/discorduser/gin-auth"
	discordcommand "backend/app/pkg/discord-command"

	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"
)

type chatController struct {
	auth ginauth.AuthInterface

	rd *domainresponse.DomainResponse

	chatUsecase domain.ChatUsecase
}

func Initialze(
	apiGroup *gin.RouterGroup, dcCmdRegister discordcommand.Register,
	auth ginauth.AuthInterface,
	rd *domainresponse.DomainResponse,
	chatUsecase domain.ChatUsecase,
) {
	ctrl := chatController{
		auth: auth,

		rd: rd,

		chatUsecase: chatUsecase,
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
