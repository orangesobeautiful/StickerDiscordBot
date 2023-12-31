package delivery

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	discordcommand "backend/app/pkg/discord-command"
	"backend/app/pkg/ginext"
	"backend/app/pkg/hserr"
	"backend/app/pkg/log"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
)

func (c *stickerController) RegisterGinRouter(e *gin.Engine) {
	e.POST("/sticker-images", ginext.BindHandlerWithStdCtx(c.AddStickerImage))
	e.GET("/stickers", ginext.BindHandlerWithStdCtx(c.ListSticker))
	e.DELETE("/stickers/:id", ginext.BindURIHandlerWithStdCtx(c.DeleteSticker))
}

func (c *stickerController) RegisterDiscordCommand(dcCmdRegister discordcommand.Register) {
	dcCmdRegister.MustAdd(
		&discordgo.ApplicationCommand{
			Name:        "sticker-add",
			Description: "新增貼圖",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "stickername",
					Description: "貼圖名稱",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "imageurl",
					Description: "貼圖網址",
					Required:    true,
				},
			},
		},
		func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			options := i.ApplicationCommandData().Options

			var req addImageReq
			for _, v := range options {
				switch v.Name {
				case "stickername":
					req.StickerName = v.StringValue()
				case "imageurl":
					req.ImageURL = v.StringValue()
				}
			}

			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			})
			if err != nil {
				slog.Error("s.InteractionRespond failed", slog.Any("err", err))
				return
			}

			_, err = c.AddStickerImage(context.Background(), &req)
			if err != nil {
				dcInteractionERROutput(s, i, err)
				return
			}

			_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Content: "新增貼圖成功",
			})
			if err != nil {
				slog.Error("s.InteractionRespond failed", slog.Any("err", err))
			}
		},
	)
}

func dcInteractionERROutput(s *discordgo.Session, i *discordgo.InteractionCreate, err error) {
	var respErr hserr.ErrResp
	if errors.As(err, &respErr) {
		statusCode := respErr.HTTPStatus()
		if statusCode >= http.StatusInternalServerError && statusCode < 600 {
			// TODO: use slog
			log.Errorf("respErr=%+v", err)
			_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Content: "內部伺服器錯誤",
			})
			if err != nil {
				slog.Error("s.FollowupMessageCreate failed", slog.Any("err", err))
			}

			return
		}

		respMessageBuilder := strings.Builder{}
		respMessageBuilder.WriteString(respErr.Message())
		for _, v := range respErr.Detail() {
			respMessageBuilder.WriteString("\n")
			respMessageBuilder.WriteString(v)
		}

		_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: respMessageBuilder.String(),
		})
		if err != nil {
			slog.Error("s.FollowupMessageCreate failed", slog.Any("err", err))
		}
		return
	}

	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: "內部伺服器錯誤",
	})
	if err != nil {
		slog.Error("s.FollowupMessageCreate failed", slog.Any("err", err))
	}
}
