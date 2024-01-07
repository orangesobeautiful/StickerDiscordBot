package delivery

import (
	discordcommand "backend/app/pkg/discord-command"
	"backend/app/pkg/ginext"

	"github.com/gin-gonic/gin"
)

func (c *discorduserController) RegisterGinRouter(apiGroup *gin.RouterGroup) {
	authRequiredMiddleware := c.auth.GetRequiredAuthMiddleware()

	apiGroup.GET("/login-code", ginext.HandlerWithStdCtx(c.CreateLoginCode))
	apiGroup.POST("/login-code", ginext.BindHandler(c.CheckLoginCode))

	apiGroup.GET("/me", authRequiredMiddleware, ginext.Handler(c.GetSelfInfo))
}

func (c *discorduserController) RegisterDiscordCommand(dcCmdRegister discordcommand.Register) {
	dcCmd, dcCmdHandler := discordcommand.DiscordCommandRegister(
		"web-login", "verify web login code",
		c.VerifyLoginCode,
	)
	dcCmdRegister.MustAdd(dcCmd, dcCmdHandler)
}
