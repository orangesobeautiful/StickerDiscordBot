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
	ragReferencePoolIDParam = "rag_ref_pool_id"
	ragReferenceTextIDParam = "rag_ref_text_id"
)

func getRAGReferencePoolIDFromContext(ctx *gin.Context) (int, error) {
	return commonDelivery.GetIDParamFromContext(ctx, ragReferencePoolIDParam)
}

func getRAGReferenceTextIDFromContext(ctx *gin.Context) (int, error) {
	return commonDelivery.GetIDParamFromContext(ctx, ragReferenceTextIDParam)
}

func (c *chatController) registerGinRouter(apiGroup *gin.RouterGroup) {
	authRequiredMiddleware := c.auth.GetRequiredAuthMiddleware()

	apiGroup.POST("/chat", authRequiredMiddleware, ginext.BindHandler(c.ginChat))

	specifyRAGReferencePoolIDGroup := apiGroup.Group("/rag_ref_pools/:" + ragReferencePoolIDParam)
	specifyRAGReferencePoolIDGroup.Use(authRequiredMiddleware, c.specifyRAGReferencePoolIDAuth)
	specifyRAGReferencePoolIDGroup.POST("/rag_ref_texts", ginext.BindHandler(c.ginCreateRAGReferencePoolTexts))
	specifyRAGReferencePoolIDGroup.GET("/rag_ref_texts", ginext.BindHandler(c.ginListRAGReferencePoolTexts))

	specifyRAGReferenceTextIDGroup := apiGroup.Group("/rag_ref_texts/:" + ragReferenceTextIDParam)
	specifyRAGReferenceTextIDGroup.Use(authRequiredMiddleware, c.specifyRAGReferenceTextIDAuth)
	specifyRAGReferenceTextIDGroup.DELETE("", ginext.Handler(c.ginDeleteRAGReferenceText))
}

func (c *chatController) specifyRAGReferencePoolIDAuth(ctx *gin.Context) {
	user := c.auth.MustGetUserFromContext(ctx)

	ragReferencePoolIDStr := ctx.Param(ragReferencePoolIDParam)
	ragReferencePoolID, err := strconv.ParseInt(ragReferencePoolIDStr, 10, 64)
	if err != nil {
		commonDelivery.GinAbortWithError(ctx,
			hserr.New(http.StatusBadRequest, "rag reference pool is not a number"))
		return
	}

	guildOwnRAGReferencePool, err := c.dcGuildUsecase.IsGuildOwnRAGReferencePool(ctx, user.GuildID, int(ragReferencePoolID))
	if err != nil {
		if domain.IsNotFoundError(err) {
			commonDelivery.GinAbortWithError(ctx,
				hserr.New(http.StatusForbidden, "you are not allowed to access this rag reference pool"))
			return
		}

		commonDelivery.GinAbortWithError(ctx,
			hserr.NewInternalError(err, "check if guild own rag reference pool"))
		return
	}

	if !guildOwnRAGReferencePool {
		commonDelivery.GinAbortWithError(ctx,
			hserr.New(http.StatusForbidden, "you are not allowed to access this rag reference pool"))
		return
	}

	ctx.Next()
}

func (c *chatController) specifyRAGReferenceTextIDAuth(ctx *gin.Context) {
	user := c.auth.MustGetUserFromContext(ctx)

	ragReferenceTextIDStr := ctx.Param(ragReferenceTextIDParam)
	ragReferenceTextID, err := strconv.ParseInt(ragReferenceTextIDStr, 10, 64)
	if err != nil {
		commonDelivery.GinAbortWithError(ctx,
			hserr.New(http.StatusBadRequest, "rag reference text is not a number"))
		return
	}

	guildOwnRAGReferenceText, err := c.dcGuildUsecase.IsGuildOwnRAGReferenceText(ctx, user.GuildID, int(ragReferenceTextID))
	if err != nil {
		if domain.IsNotFoundError(err) {
			commonDelivery.GinAbortWithError(ctx,
				hserr.New(http.StatusForbidden, "you are not allowed to access this rag reference text"))
			return
		}

		commonDelivery.GinAbortWithError(ctx,
			hserr.NewInternalError(err, "check if guild own rag reference text"))
		return
	}

	if !guildOwnRAGReferenceText {
		commonDelivery.GinAbortWithError(ctx,
			hserr.New(http.StatusForbidden, "you are not allowed to access this rag reference text"))
		return
	}

	ctx.Next()
}

func (c *chatController) registerDiscordCommand(dcCmdRegister discordcommand.Register) {
	dcCmd, dcCmdHandler := discordcommand.DiscordCommandRegister(
		"chat", "chat with bot",
		c.discordChat,
	)
	dcCmdRegister.MustAdd(dcCmd, dcCmdHandler)
}
