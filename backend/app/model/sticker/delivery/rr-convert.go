package delivery

import (
	"fmt"

	"backend/app/domain"
	domainresponse "backend/app/domain-response"
	discordcommand "backend/app/pkg/discord-command"

	"github.com/bwmarrin/discordgo"
)

type addImageReq struct {
	StickerName string `json:"sticker_name" binding:"required"`

	ImageURL string `json:"image_url" binding:"required,http_url"`
}

type listStickerReq struct {
	Page int `form:"page" binding:"required,gte=1"`

	Limit int `form:"limit" binding:"required,gte=1,lte=30"`

	Search string `form:"search"`
}

var _ discordcommand.DiscordWebhookParamsMarshaler = (*listStickerResp)(nil)

type listStickerResp struct {
	TotalCount int `json:"total_count"`

	Stickers []*domainresponse.Sticker `json:"stickers"`
}

func (c *stickerController) newlistStickerRespFromListResult(listResult domain.ListStickerResult) *listStickerResp {
	entSs := listResult.GetItems()
	ss := make([]*domainresponse.Sticker, len(entSs))
	for i, entS := range entSs {
		ss[i] = c.rd.NewStickerFromEnt(entS)
	}

	return &listStickerResp{
		TotalCount: listResult.GetTotal(),
		Stickers:   ss,
	}
}

func (l *listStickerResp) MarshalDiscordWebhookParams() *discordgo.WebhookParams {
	result := new(discordgo.WebhookParams)

	result.Content = "貼圖列表"

	for _, s := range l.Stickers {
		stickerEmbedTitle := fmt.Sprintf("%d: %s", s.ID, s.StickerName)

		for _, img := range s.Images {
			imgEmbed := &discordgo.MessageEmbed{
				Type:  discordgo.EmbedTypeImage,
				Title: stickerEmbedTitle,
				URL:   fmt.Sprintf("https://%d.for.display.multiple.images.in.a.single.embed", s.ID),
				Image: &discordgo.MessageEmbedImage{
					URL: img.URL,
				},
			}

			result.Embeds = append(result.Embeds, imgEmbed)
		}
	}

	return result
}

type deleteStickerReq struct {
	ID int `uri:"id" binding:"required,gte=0"`
}

type deleteStickerByNameReq struct {
	Name string `json:"name" binding:"required"`
}
