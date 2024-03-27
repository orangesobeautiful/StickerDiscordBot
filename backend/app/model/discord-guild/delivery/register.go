package delivery

import (
	"net/http"
	"strconv"

	"backend/app/domain"
	commonDelivery "backend/app/pkg/common/delivery"
	discordcommand "backend/app/pkg/discord-command"
	"backend/app/pkg/ginext"
	"backend/app/pkg/hserr"

	"github.com/gin-gonic/gin"
)

const (
	guildIDParam    = "guild_id"
	chatroomIDParam = "chatroom_id"
)

func (c *discordGuildController) registerGinRouter(apiGroup *gin.RouterGroup) {
	authRequiredMiddleware := c.auth.GetRequiredAuthMiddleware()

	specifyGuildIDGroup := apiGroup.Group("/guilds/:" + guildIDParam)
	specifyGuildIDGroup.Use(authRequiredMiddleware)
	specifyGuildIDGroup.Use(c.specifyGuildIDAuth)

	specifyGuildIDGroup.POST("/chatrooms", ginext.BindHandler(c.ginCreateGuildChatroom))
	specifyGuildIDGroup.GET("/chatrooms", ginext.BindHandler(c.ginlistGuildChatrooms))

	specifyChatroomIDGroup := apiGroup.Group("/chatrooms/:" + chatroomIDParam)
	specifyChatroomIDGroup.Use(c.specifyChatroomIDAuth)

	specifyChatroomIDGroup.DELETE("", ginext.BindURIHandler(c.ginDeleteGuildChatroom))
}

func (c *discordGuildController) registerDiscordCommand(dcCmdRegister discordcommand.Register) {
	dcCmd, dcCmdHandler := discordcommand.DiscordCommandRegister(
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
