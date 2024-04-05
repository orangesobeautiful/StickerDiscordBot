package delivery

import (
	"context"
	"net/http"

	"backend/app/domain"
	commonDelivery "backend/app/pkg/common/delivery"
	discordcommand "backend/app/pkg/discord-command"
	"backend/app/pkg/ginext"
	"backend/app/pkg/hserr"

	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"
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

func (c *chatController) isGinUserGuildOwnResource(
	ctx *gin.Context, resourceIDParam, resourceLabel string,
	guildOwnResourceFunc func(ctx context.Context, guildID string, ragReferencePoolID int) (isOwn bool, err error),
) {
	user := c.auth.MustGetUserFromContext(ctx)

	resourceID, err := commonDelivery.GetIDParamFromContext(ctx, resourceIDParam)
	if err != nil {
		err = xerrors.Errorf("get %s id from context: %w", resourceLabel, err)
		commonDelivery.GinAbortWithError(ctx, err)
	}

	isGuildOwnResource, err := guildOwnResourceFunc(ctx, user.GuildID, resourceID)
	if err != nil {
		if domain.IsNotFoundError(err) {
			commonDelivery.GinAbortWithError(ctx, hserr.New(http.StatusForbidden, "you are not allowed to access this "+resourceLabel))
			return
		}

		commonDelivery.GinAbortWithError(ctx, hserr.NewInternalError(err, "check if guild own "+resourceLabel))
		return
	}

	if !isGuildOwnResource {
		commonDelivery.GinAbortWithError(ctx, hserr.New(http.StatusForbidden, "you are not allowed to access this "+resourceLabel))
		return
	}
}

func (c *chatController) specifyRAGReferencePoolIDAuth(ctx *gin.Context) {
	c.isGinUserGuildOwnResource(ctx, ragReferencePoolIDParam, "rag reference pool", c.dcGuildUsecase.IsGuildOwnRAGReferencePool)
}

func (c *chatController) specifyRAGReferenceTextIDAuth(ctx *gin.Context) {
	c.isGinUserGuildOwnResource(ctx, ragReferenceTextIDParam, "rag reference text", c.dcGuildUsecase.IsGuildOwnRAGReferenceText)
}

func (c *chatController) registerDiscordCommand(dcCmdRegister discordcommand.Register) {
	dcCmd, dcCmdHandler := discordcommand.DiscordCommandRegister(
		"chat", "chat with bot",
		c.discordChat,
	)
	dcCmdRegister.MustAdd(dcCmd, dcCmdHandler)
}
