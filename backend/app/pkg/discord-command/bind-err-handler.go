package discordcommand

import (
	"errors"
	"log/slog"
	"net/http"

	"backend/app/pkg/hserr"

	"github.com/bwmarrin/discordgo"
)

var ValidateErrorConvert func(err error) error

func dcInteractionErrResponse(s *discordgo.Session, i *discordgo.InteractionCreate, err error) {
	var respErr hserr.ErrResp
	if errors.As(err, &respErr) {
		statusCode := respErr.HTTPStatus()
		if statusCode >= http.StatusInternalServerError && statusCode < 600 {
			internalServerErrorResp(s, i, err)
			return
		}

		webhookParams := convertHsErrToWebhookParams(respErr)
		_, err = s.FollowupMessageCreate(i.Interaction, true, webhookParams)
		if err != nil {
			slog.Error("s.FollowupMessageCreate failed", slog.Any("err", err))
		}
		return
	}

	internalServerErrorResp(s, i, err)
}

func internalServerErrorResp(s *discordgo.Session, i *discordgo.InteractionCreate, err error) {
	slog.Error("error occurred", slog.Any("err", err))

	_, err = s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: "internal server error",
	})
	if err != nil {
		slog.Error("s.FollowupMessageCreate failed", slog.Any("err", err))
	}
}

func convertHsErrToWebhookParams(respErr hserr.ErrResp) *discordgo.WebhookParams {
	webhookParams := &discordgo.WebhookParams{
		Content: respErr.GetMessage(),
	}
	embeds := make([]*discordgo.MessageEmbed, 0, len(respErr.GetDetails()))
	for _, v := range respErr.GetDetails() {
		embeds = append(embeds, &discordgo.MessageEmbed{
			Title: v,
		})
	}
	webhookParams.Embeds = embeds
	return webhookParams
}
