package delivery

import (
	"context"

	discordcommand "backend/app/pkg/discord-command"
	"backend/app/pkg/ginext"

	"github.com/gin-gonic/gin"
)

func (c *stickerController) RegisterGinRouter(e *gin.Engine) {
	e.POST("/sticker-images", ginext.BindHandlerWithStdCtx(c.AddStickerImage))
	e.GET("/stickers", ginext.BindHandlerWithStdCtx(c.ListSticker))
	e.DELETE("/stickers/:id", ginext.BindURIHandlerWithStdCtx(c.DeleteSticker))
}

func (c *stickerController) RegisterDiscordCommand(dcCmdRegister discordcommand.Register) {
	dcCmd, dcCmdHandler := discordcommand.DiscordCommandRegister(
		"sticker-add", "新增貼圖",
		func(req *addImageReq) (*ginext.EmptyResp, error) {
			return c.AddStickerImage(context.Background(), req)
		},
	)
	dcCmdRegister.MustAdd(dcCmd, dcCmdHandler)
}
