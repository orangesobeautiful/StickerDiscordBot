package delivery

import (
	discordcommand "backend/app/pkg/discord-command"
	"backend/app/pkg/ginext"

	"github.com/gin-gonic/gin"
)

func (c *chatController) registerGinRouter(apiGroup *gin.RouterGroup) {
	authRequiredMiddleware := c.auth.GetRequiredAuthMiddleware()

	apiGroup.POST("/chat", authRequiredMiddleware, ginext.BindHandler(c.ginChat))
}

func (c *chatController) registerDiscordCommand(dcCmdRegister discordcommand.Register) {
	dcCmd, dcCmdHandler := discordcommand.DiscordCommandRegister(
		"chat", "chat with bot",
		c.discordChat,
	)
	dcCmdRegister.MustAdd(dcCmd, dcCmdHandler)
}
