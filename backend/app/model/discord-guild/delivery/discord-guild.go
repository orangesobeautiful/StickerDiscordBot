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

type discordGuildController struct {
	auth ginauth.AuthInterface

	rd *domainresponse.DomainResponse

	dcGuildUsecase domain.DiscordGuildUsecase
}

func Initialze(
	apiGroup *gin.RouterGroup, dcCmdRegister discordcommand.Register,
	auth ginauth.AuthInterface,
	rd *domainresponse.DomainResponse,
	dcGuildUsecase domain.DiscordGuildUsecase,
) {
	ctrl := &discordGuildController{
		auth: auth,

		rd: rd,

		dcGuildUsecase: dcGuildUsecase,
	}

	ctrl.registerGinRouter(apiGroup)
	ctrl.registerDiscordCommand(dcCmdRegister)
}

func (c *discordGuildController) ginCreateGuildChatroom(
	ctx *gin.Context, req *ginCreateGuildChatroomReq,
) (resp *createGuildChatroomResp, err error) {
	user := c.auth.MustGetUserFromContext(ctx)

	chatroomID, err := c.dcGuildUsecase.CreateGuildChatroom(ctx, user.GuildID, req.Name)
	if err != nil {
		return nil, xerrors.Errorf("create guild chatroom: %w", err)
	}

	return &createGuildChatroomResp{
		ChatroomID: chatroomID,
	}, nil
}

func (c *discordGuildController) discordCreateGuildChatroom(
	req *discordCreateGuildChatroomReq,
) (resp *createGuildChatroomResp, err error) {
	ctx := context.Background()
	chatroomID, err := c.dcGuildUsecase.CreateGuildChatroom(ctx, req.GuildID, req.Name)
	if err != nil {
		return nil, xerrors.Errorf("create guild chatroom: %w", err)
	}

	return &createGuildChatroomResp{
		ChatroomID: chatroomID,
	}, nil
}

func (c *discordGuildController) ginlistGuildChatrooms(
	ctx *gin.Context, req *ginlistGuildChatroomsReq,
) (resp *listChatroomsResp, err error) {
	user := c.auth.MustGetUserFromContext(ctx)

	offset := (req.Page - 1) * req.Limit

	chatroomsResult, err := c.dcGuildUsecase.ListGuildChatrooms(ctx, user.GuildID, req.Limit, offset)
	if err != nil {
		return nil, xerrors.Errorf("list guild chatrooms: %w", err)
	}

	return newlistChatroomRespFromListResult(chatroomsResult), nil
}

const discordListGuildChatroomLimit = 10

func (c *discordGuildController) discordListGuildChatrooms(
	req *discordListGuildChatroomsReq,
) (resp *listChatroomsResp, err error) {
	ctx := context.Background()

	limit := discordListGuildChatroomLimit
	offset := (req.Page - 1) * limit

	chatroomsResult, err := c.dcGuildUsecase.ListGuildChatrooms(ctx, req.GuildID, limit, offset)
	if err != nil {
		return nil, xerrors.Errorf("list guild chatrooms: %w", err)
	}

	return newlistChatroomRespFromListResult(chatroomsResult), nil
}

func (c *discordGuildController) ginDeleteGuildChatroom(
	ctx *gin.Context, req *ginDeleteGuildChatroomReq,
) (resp *ginext.EmptyResp, err error) {
	err = c.dcGuildUsecase.RemoveGuildChatroom(ctx, req.ChatroomID)
	if err != nil {
		return nil, xerrors.Errorf("delete guild chatroom: %w", err)
	}

	return &ginext.EmptyResp{}, nil
}

func (c *discordGuildController) discordDeleteGuildChatroom(
	req *discordDeleteGuildChatroomReq,
) (resp *ginext.EmptyResp, err error) {
	ctx := context.Background()
	err = c.dcGuildUsecase.RemoveGuildChatroom(ctx, req.ChatroomID)
	if err != nil {
		return nil, xerrors.Errorf("delete guild chatroom: %w", err)
	}

	return &ginext.EmptyResp{}, nil
}
