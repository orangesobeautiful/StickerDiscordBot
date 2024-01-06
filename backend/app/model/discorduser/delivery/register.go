package delivery

import (
	discordcommand "backend/app/pkg/discord-command"
	"backend/app/pkg/ginext"

	"github.com/gin-gonic/gin"
)

func (c *discorduserController) RegisterGinRouter(e *gin.Engine) {
	authRequiredMiddleware := c.auth.GetRequiredAuthMiddleware()

	e.GET("/login-code", ginext.HandlerWithStdCtx(c.CreateLoginCode))
	e.POST("/login-code", ginext.BindHandler(c.CheckLoginCode))

	e.GET("/me", authRequiredMiddleware, ginext.Handler(c.GetSelfInfo))
}

func (c *discorduserController) RegisterDiscordCommand(dcCmdRegister discordcommand.Register) {
	dcCmd, dcCmdHandler := discordcommand.DiscordCommandRegister(
		"web-login", "verify web login code",
		c.VerifyLoginCode,
	)
	dcCmdRegister.MustAdd(dcCmd, dcCmdHandler)
}
