package delivery

import (
	"context"
	"net/http"
	"strconv"

	"backend/app/domain"
	commonDelivery "backend/app/pkg/common/delivery"
	discordcommand "backend/app/pkg/discord-command"
	"backend/app/pkg/ginext"
	"backend/app/pkg/hserr"

	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"
)

const (
	guildIDParam    = "guild_id"
	stickerIDParam  = "sticker_id"
	chatroomIDParam = "chatroom_id"
)

func getGuildIDFromContext(ctx *gin.Context) string {
	return ctx.Param(guildIDParam)
}

func getStickerIDFromContext(ctx *gin.Context) (int, error) {
	stickerIDStr := ctx.Param(stickerIDParam)
	stickerID, err := strconv.Atoi(stickerIDStr)
	if err != nil {
		return 0, hserr.New(http.StatusBadRequest, "sticker id is not a number")
	}

	return stickerID, nil
}

func getChatroomIDFromContext(ctx *gin.Context) (int, error) {
	chatroomIDStr := ctx.Param(chatroomIDParam)
	chatroomID, err := strconv.Atoi(chatroomIDStr)
	if err != nil {
		return 0, hserr.New(http.StatusBadRequest, "chatroom id is not a number")
	}
	return chatroomID, nil
}

func (c *discordGuildController) registerGinRouter(apiGroup *gin.RouterGroup) {
	authRequiredMiddleware := c.auth.GetRequiredAuthMiddleware()

	specifyGuildIDGroup := apiGroup.Group("/guilds/:" + guildIDParam)
	specifyGuildIDGroup.Use(authRequiredMiddleware)
	specifyGuildIDGroup.Use(c.specifyGuildIDAuth)

	specifyGuildIDGroup.POST("/sticker-images", ginext.BindHandler(c.ginAddStickerImage))
	specifyGuildIDGroup.GET("/stickers", ginext.BindHandler(c.ginListSticker))
	specifyGuildIDGroup.GET("/sticker_by_name", ginext.BindHandler(c.ginGetStickerByName))

	specifyStickerIDGroup := apiGroup.Group("/stickers/:" + stickerIDParam)
	specifyStickerIDGroup.Use(authRequiredMiddleware, c.specifyStickerIDAuth)
	specifyStickerIDGroup.DELETE("", ginext.BindURIHandler(c.ginDeleteSticker))

	specifyGuildIDGroup.POST("/chatrooms", ginext.BindHandler(c.ginCreateGuildChatroom))
	specifyGuildIDGroup.GET("/chatrooms", ginext.BindHandler(c.ginlistGuildChatrooms))
	specifyGuildIDGroup.POST("/rag_ref_pools", ginext.BindHandler(c.ginCreateGuildRAGReferencePool))
	specifyGuildIDGroup.GET("/rag_ref_pools", ginext.BindHandler(c.ginListGuildRAGReferencePools))

	specifyChatroomIDGroup := apiGroup.Group("/chatrooms/:" + chatroomIDParam)
	specifyChatroomIDGroup.Use(authRequiredMiddleware, c.specifyChatroomIDAuth)

	specifyChatroomIDGroup.DELETE("", ginext.BindURIHandler(c.ginDeleteGuildChatroom))
	specifyChatroomIDGroup.POST("/rag_ref_pools", ginext.BindHandler(c.ginAddChatroomRAGReferencePool))
	specifyChatroomIDGroup.GET("/rag_ref_pools", ginext.BindHandler(c.ginListChatroomRAGReferencePools))
	specifyChatroomIDGroup.DELETE("/rag_ref_pools", ginext.BindHandler(c.ginRemoveChatroomRAGReferencePools))
}

func (c *discordGuildController) registerDiscordCommand(dcCmdRegister discordcommand.Register) {
	dcCmd, dcCmdHandler := discordcommand.DiscordCommandRegister(
		"sticker-add", "新增貼圖",
		func(req *discordAddImageReq) (*ginext.EmptyResp, error) {
			err := c.addStickerImage(context.Background(), req.GuildID, req.StickerName, req.ImageURL)
			if err != nil {
				return nil, xerrors.Errorf("add sticker image: %w", err)
			}

			return nil, nil
		},
	)
	dcCmdRegister.MustAdd(dcCmd, dcCmdHandler)

	dcCmd, dcCmdHandler = discordcommand.DiscordCommandRegister(
		"sticker-list", "列出貼圖",
		func(req *discordListStickerReq) (*listStickerResp, error) {
			return c.listSticker(context.Background(), req.GuildID, req.Page, req.Limit, req.Search)
		},
	)
	dcCmdRegister.MustAdd(dcCmd, dcCmdHandler)

	dcCmd, dcCmdHandler = discordcommand.DiscordCommandRegister(
		"sticker-delete", "刪除貼圖",
		func(req *discordDeleteStickerReq) (*ginext.EmptyResp, error) {
			ctx := context.Background()

			guildOwn, err := c.dcGuildUsecase.IsGuildOwnSticker(ctx, req.GuildID, req.ID)
			if err != nil {
				return nil, xerrors.Errorf("check if guild own sticker: %w", err)
			}
			if !guildOwn {
				return nil, hserr.New(http.StatusForbidden, "you are not allowed to access this sticker")
			}

			err = c.deleteSticker(context.Background(), req.ID)
			if err != nil {
				return nil, xerrors.Errorf("delete sticker: %w", err)
			}

			return nil, nil
		},
	)
	dcCmdRegister.MustAdd(dcCmd, dcCmdHandler)

	dcCmd, dcCmdHandler = discordcommand.DiscordCommandRegister(
		"sticker-delete-by-name", "刪除貼圖",
		func(req *discordDeleteStickerByNameReq) (*ginext.EmptyResp, error) {
			err := c.deleteStickerByName(context.Background(), req.GuildID, req.Name)
			if err != nil {
				return nil, xerrors.Errorf("delete sticker by name: %w", err)
			}

			return nil, nil
		},
	)
	dcCmdRegister.MustAdd(dcCmd, dcCmdHandler)

	dcCmd, dcCmdHandler = discordcommand.DiscordCommandRegister(
		"create_chatroom", "create new chatroom in guild",
		c.discordCreateGuildChatroom,
	)
	dcCmdRegister.MustAdd(dcCmd, dcCmdHandler)

	dcCmd, dcCmdHandler = discordcommand.DiscordCommandRegister(
		"list_chatrooms", "list chatrooms in guild",
		c.discordListGuildChatrooms,
	)
	dcCmdRegister.MustAdd(dcCmd, dcCmdHandler)

	dcCmd, dcCmdHandler = discordcommand.DiscordCommandRegister(
		"delete_chatroom", "delete chatroom in guild",
		c.discordDeleteGuildChatroom,
	)
	dcCmdRegister.MustAdd(dcCmd, dcCmdHandler)
}

func (c *discordGuildController) specifyGuildIDAuth(ctx *gin.Context) {
	user := c.auth.MustGetUserFromContext(ctx)

	guildID := ctx.Param(guildIDParam)
	if guildID != user.GuildID {
		commonDelivery.GinAbortWithError(ctx,
			hserr.New(http.StatusForbidden, "you are not allowed to access this guild"))
		return
	}

	ctx.Next()
}

func (c discordGuildController) specifyStickerIDAuth(ctx *gin.Context) {
	user := c.auth.MustGetUserFromContext(ctx)

	stickerID, err := getStickerIDFromContext(ctx)
	if err != nil {
		commonDelivery.GinAbortWithError(ctx, err)
		return
	}

	guildOwnSticker, err := c.dcGuildUsecase.IsGuildOwnSticker(ctx, user.GuildID, stickerID)
	if err != nil {
		if domain.IsNotFoundError(err) {
			commonDelivery.GinAbortWithError(ctx,
				hserr.New(http.StatusForbidden, "you are not allowed to access this sticker"))
			return
		}

		commonDelivery.GinAbortWithError(ctx,
			hserr.NewInternalError(err, "check if guild own sticker"))
		return
	}

	if !guildOwnSticker {
		commonDelivery.GinAbortWithError(ctx,
			hserr.New(http.StatusForbidden, "you are not allowed to access this sticker"))
		return
	}

	ctx.Next()
}

func (c *discordGuildController) specifyChatroomIDAuth(ctx *gin.Context) {
	user := c.auth.MustGetUserFromContext(ctx)

	chatroomIDStr := ctx.Param(chatroomIDParam)
	chatroomID, err := strconv.ParseInt(chatroomIDStr, 10, 64)
	if err != nil {
		commonDelivery.GinAbortWithError(ctx,
			hserr.New(http.StatusBadRequest, "chatroom id is not a number"))
		return
	}

	guildOwnChatroom, err := c.dcGuildUsecase.IsGuildOwnChatroom(ctx, user.GuildID, int(chatroomID))
	if err != nil {
		if domain.IsNotFoundError(err) {
			commonDelivery.GinAbortWithError(ctx,
				hserr.New(http.StatusForbidden, "you are not allowed to access this chatroom"))
			return
		}

		commonDelivery.GinAbortWithError(ctx,
			hserr.NewInternalError(err, "check if guild own chatroom"))
		return
	}

	if !guildOwnChatroom {
		commonDelivery.GinAbortWithError(ctx,
			hserr.New(http.StatusForbidden, "you are not allowed to access this chatroom"))
		return
	}

	ctx.Next()
}
