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
	guildIDParam        = "guild_id"
	stickerIDParam      = "sticker_id"
	stickerImageIDParam = "sticker_image_id"
	chatroomIDParam     = "chatroom_id"
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

func getStickerImageIDFromContext(ctx *gin.Context) (int, error) {
	stickerImageIDStr := ctx.Param(stickerImageIDParam)
	stickerImageID, err := strconv.Atoi(stickerImageIDStr)
	if err != nil {
		return 0, hserr.New(http.StatusBadRequest, "sticker image id is not a number")
	}

	return stickerImageID, nil
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

	specifyStickerImageIDGroup := apiGroup.Group("/sticker-images/:" + stickerImageIDParam)
	specifyStickerImageIDGroup.Use(authRequiredMiddleware, c.specifyStickerImageIDAuth)
	specifyStickerImageIDGroup.DELETE("", ginext.BindURIHandler(c.ginDeleteStickerImage))

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
		"sticker-image-delete", "刪除貼圖圖片",
		func(req *discordDeleteStickerImageReq) (*ginext.EmptyResp, error) {
			ctx := context.Background()

			guildOwn, err := c.dcGuildUsecase.IsGuildOwnStickerImage(ctx, req.GuildID, req.ID)
			if err != nil {
				return nil, xerrors.Errorf("check if guild own sticker image: %w", err)
			}

			if !guildOwn {
				return nil, hserr.New(http.StatusForbidden, "you are not allowed to access this sticker")
			}

			err = c.deleteStickerImage(ctx, req.ID)
			if err != nil {
				return nil, xerrors.Errorf("delete sticker image: %w", err)
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
	err := c.verifyUserGuildOwnResource(
		ctx,
		"sticker",
		getStickerIDFromContext,
		c.dcGuildUsecase.IsGuildOwnSticker,
	)
	if err != nil {
		commonDelivery.GinAbortWithError(ctx, err)

		return
	}

	ctx.Next()
}

func (c discordGuildController) specifyStickerImageIDAuth(ctx *gin.Context) {
	err := c.verifyUserGuildOwnResource(
		ctx,
		"sticker image",
		getStickerImageIDFromContext,
		c.dcGuildUsecase.IsGuildOwnStickerImage,
	)
	if err != nil {
		commonDelivery.GinAbortWithError(ctx, err)

		return
	}

	ctx.Next()
}

func (c *discordGuildController) specifyChatroomIDAuth(ctx *gin.Context) {
	err := c.verifyUserGuildOwnResource(
		ctx,
		"chatroom",
		getChatroomIDFromContext,
		c.dcGuildUsecase.IsGuildOwnChatroom,
	)
	if err != nil {
		commonDelivery.GinAbortWithError(ctx, err)

		return
	}

	ctx.Next()
}

func (c *discordGuildController) verifyUserGuildOwnResource(
	ctx *gin.Context,
	resourceName string,
	getResourceIDFunc func(ctx *gin.Context) (int, error),
	isOwnResourceFunc func(ctx context.Context, guildID string, resourceID int) (bool, error),
) error {
	user := c.auth.MustGetUserFromContext(ctx)

	resourceID, err := getResourceIDFunc(ctx)
	if err != nil {
		return err
	}

	isOwn, err := isOwnResourceFunc(ctx, user.GuildID, resourceID)
	if err != nil {
		if domain.IsNotFoundError(err) {
			return hserr.New(http.StatusForbidden, "you are not allowed to access this "+resourceName)
		}

		return xerrors.Errorf("check if guild own resource: %w", err)
	}

	if !isOwn {
		return hserr.New(http.StatusForbidden, "you are not allowed to access this "+resourceName)
	}

	return nil
}
