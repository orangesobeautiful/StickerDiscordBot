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

	stickerUsecase domain.StickerUsecase
}

func Initialze(
	apiGroup *gin.RouterGroup, dcCmdRegister discordcommand.Register,
	auth ginauth.AuthInterface,
	rd *domainresponse.DomainResponse,
	dcGuildUsecase domain.DiscordGuildUsecase, stickerUsecase domain.StickerUsecase,
) {
	ctrl := &discordGuildController{
		auth: auth,

		rd: rd,

		dcGuildUsecase: dcGuildUsecase,

		stickerUsecase: stickerUsecase,
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
	ctx *gin.Context, req ginlistGuildChatroomsReq,
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

func (c *discordGuildController) ginCreateGuildRAGReferencePool(
	ctx *gin.Context, req *createGuildRAGReferencePoolReq,
) (resp *createGuildRAGReferencePoolResp, err error) {
	guildID := getGuildIDFromContext(ctx)

	return c.createGuildRAGReferencePool(ctx, guildID, req)
}

func (c *discordGuildController) createGuildRAGReferencePool(
	ctx *gin.Context, guildID string, req *createGuildRAGReferencePoolReq,
) (resp *createGuildRAGReferencePoolResp, err error) {
	id, err := c.dcGuildUsecase.CreateRAGReferencePool(ctx, guildID, req.Name, req.Description)
	if err != nil {
		return nil, xerrors.Errorf("create rag reference pool: %w", err)
	}

	return &createGuildRAGReferencePoolResp{
		ID: id,
	}, nil
}

func (c *discordGuildController) ginListGuildRAGReferencePools(
	ctx *gin.Context, req listGuildRAGReferencePoolsReq,
) (resp *listGuildRAGReferencePoolsResp, err error) {
	guildID := getGuildIDFromContext(ctx)

	return c.listGuildRAGReferencePools(ctx, guildID, &req)
}

func (c *discordGuildController) listGuildRAGReferencePools(
	ctx *gin.Context, guildID string, req *listGuildRAGReferencePoolsReq,
) (resp *listGuildRAGReferencePoolsResp, err error) {
	offset := (req.Page - 1) * req.Limit

	result, err := c.dcGuildUsecase.ListRAGReferencePools(ctx, guildID, req.Limit, offset)
	if err != nil {
		return nil, xerrors.Errorf("list rag reference pools: %w", err)
	}

	return newlistGuildRAGReferencePoolsRespFromListResult(result), nil
}

func (c *discordGuildController) ginAddChatroomRAGReferencePool(
	ctx *gin.Context, req *ginAddChatroomRAGReferencePoolReq,
) (resp *ginext.EmptyResp, err error) {
	chatroomID, err := getChatroomIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return c.addChatroomRAGReferencePool(ctx, chatroomID, req.RAGReferencePoolID)
}

func (c *discordGuildController) addChatroomRAGReferencePool(
	ctx *gin.Context, chatroomID, ragReferencePoolID int,
) (resp *ginext.EmptyResp, err error) {
	err = c.dcGuildUsecase.AddChatroomRAGReferencePool(ctx, chatroomID, ragReferencePoolID)
	if err != nil {
		return nil, xerrors.Errorf("add chatroom rag reference pool: %w", err)
	}

	return &ginext.EmptyResp{}, nil
}

func (c *discordGuildController) ginListChatroomRAGReferencePools(
	ctx *gin.Context, req *ginListChatroomRAGReferencePoolsReq,
) (resp *listChatroomRAGReferencePoolsResp, err error) {
	chatroomID, err := getChatroomIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return c.listChatroomRAGReferencePools(ctx, chatroomID, req)
}

func (c *discordGuildController) listChatroomRAGReferencePools(
	ctx *gin.Context, chatroomID int, req *ginListChatroomRAGReferencePoolsReq,
) (resp *listChatroomRAGReferencePoolsResp, err error) {
	offset := (req.Page - 1) * req.Limit

	result, err := c.dcGuildUsecase.ListChatroomRAGReferencePools(ctx, chatroomID, req.Limit, offset)
	if err != nil {
		return nil, xerrors.Errorf("list chatroom rag reference pools: %w", err)
	}

	return newlistChatroomRAGReferencePoolsRespFromListResult(result), nil
}

func (c *discordGuildController) ginRemoveChatroomRAGReferencePools(
	ctx *gin.Context, req *ginRemoveChatroomRAGReferencePoolsReq,
) (resp *ginext.EmptyResp, err error) {
	chatroomID, err := getChatroomIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return c.removeChatroomRAGReferencePools(ctx, chatroomID, req.RAGReferencePoolIDs)
}

func (c *discordGuildController) removeChatroomRAGReferencePools(
	ctx *gin.Context, chatroomID int, ragReferencePoolIDs []int,
) (resp *ginext.EmptyResp, err error) {
	err = c.dcGuildUsecase.RemoveChatroomRAGReferencePools(ctx, chatroomID, ragReferencePoolIDs)
	if err != nil {
		return nil, xerrors.Errorf("remove all chatroom rag reference pools: %w", err)
	}

	return &ginext.EmptyResp{}, nil
}
