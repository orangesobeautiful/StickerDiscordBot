package delivery

import (
	"fmt"
	"time"

	"backend/app/domain"
	domainresponse "backend/app/domain-response"
	"backend/app/ent"
	discordcommand "backend/app/pkg/discord-command"

	"github.com/bwmarrin/discordgo"
)

type ginAddImageReq struct {
	StickerName string `json:"sticker_name" binding:"required"`

	ImageURL string `json:"image_url" binding:"required,http_url"`
}

type discordAddImageReq struct {
	discordcommand.BaseAuthInteractionCreate `dccmd:"ignore"`

	StickerName string `dccmd:"name=sticker_name" binding:"required"`

	ImageURL string `dccmd:"name=image_url" binding:"required,http_url"`
}

type ginListStickerReq struct {
	Page int `form:"page" binding:"required,gte=1"`

	Limit int `form:"limit" binding:"required,gte=1,lte=30"`

	Search string `form:"search"`
}

type discordListStickerReq struct {
	discordcommand.BaseAuthInteractionCreate `dccmd:"ignore"`

	Page int `dccmd:"name=page" binding:"required,gte=1"`

	Limit int `dccmd:"name=limit" binding:"required,gte=1,lte=30"`

	Search string `dccmd:"name=search"`
}

var _ discordcommand.DiscordWebhookParamsMarshaler = (*listStickerResp)(nil)

type listStickerResp struct {
	TotalCount int `json:"total_count"`

	Stickers []*domainresponse.Sticker `json:"stickers"`
}

func (c *discordGuildController) newlistStickerRespFromListResult(listResult domain.ListStickerResult) *listStickerResp {
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

type ginDeleteStickerReq struct {
	StickerID int `uri:"sticker_id" binding:"required,gte=0"`
}

type discordDeleteStickerReq struct {
	discordcommand.BaseAuthInteractionCreate `dccmd:"ignore"`

	ID int `dccmd:"name=id" binding:"required,gte=0"`
}

type discordDeleteStickerByNameReq struct {
	discordcommand.BaseAuthInteractionCreate `dccmd:"ignore"`

	Name string `dccmd:"name=name" binding:"required"`
}

type ginCreateGuildChatroomReq struct {
	Name string `json:"name" binding:"required"`
}

type discordCreateGuildChatroomReq struct {
	discordcommand.BaseAuthInteractionCreate `dccmd:"ignore"`

	Name string `dccmd:"name=name" binding:"required"`
}

type createGuildChatroomResp struct {
	ChatroomID int `json:"chatroom_id"`
}

type ginlistGuildChatroomsReq struct {
	Page int `form:"page" binding:"required,gte=1"`

	Limit int `form:"limit" binding:"required,gte=1,lte=30"`
}

type discordListGuildChatroomsReq struct {
	discordcommand.BaseAuthInteractionCreate `dccmd:"ignore"`

	Page int `dccmd:"name=page" binding:"required,gte=1"`
}

var _ discordcommand.DiscordWebhookParamsMarshaler = (*listChatroomsResp)(nil)

type listChatroomsResp struct {
	TotalCount int `json:"total_count"`

	Chatroom []*Chatroom `json:"chatroom"`
}

type Chatroom struct {
	ID int `json:"id"`

	Name string `json:"name"`

	CreatedAt time.Time `json:"created_at"`
}

func newlistChatroomRespFromListResult(listResult domain.ListChatroomsResult) *listChatroomsResp {
	entChatrooms := listResult.GetItems()
	chatrooms := make([]*Chatroom, len(entChatrooms))
	for i, entChatroom := range entChatrooms {
		chatrooms[i] = newChatroomFromEnt(entChatroom)
	}

	return &listChatroomsResp{
		TotalCount: listResult.GetTotal(),
		Chatroom:   chatrooms,
	}
}

func newChatroomFromEnt(entChatroom *ent.Chatroom) *Chatroom {
	return &Chatroom{
		ID:        entChatroom.ID,
		Name:      entChatroom.Name,
		CreatedAt: entChatroom.CreatedAt,
	}
}

func (r *listChatroomsResp) MarshalDiscordWebhookParams() *discordgo.WebhookParams {
	result := new(discordgo.WebhookParams)

	result.Content = fmt.Sprintf("聊天室列表 總共: %d 最後一頁:%d\n",
		r.TotalCount, r.TotalCount/discordListGuildChatroomLimit+1)

	for _, chatroom := range r.Chatroom {
		chatroomEmbedTitle := fmt.Sprintf("%d: %s", chatroom.ID, chatroom.Name)

		chatroomEmbed := &discordgo.MessageEmbed{
			Title: chatroomEmbedTitle,
		}

		result.Embeds = append(result.Embeds, chatroomEmbed)
	}

	return result
}

type ginDeleteGuildChatroomReq struct {
	ChatroomID int `uri:"chatroom_id" binding:"required"`
}

type discordDeleteGuildChatroomReq struct {
	discordcommand.BaseAuthInteractionCreate `dccmd:"ignore"`

	ChatroomID int `dccmd:"name=chatroom_id" binding:"required"`
}

type createGuildRAGReferencePoolReq struct {
	Name string `json:"name" binding:"required"`

	Description string `json:"description"`
}

type createGuildRAGReferencePoolResp struct {
	ID int `json:"id"`
}

type listGuildRAGReferencePoolsReq struct {
	Page int `form:"page" binding:"required,gte=1"`

	Limit int `form:"limit" binding:"required,gte=1,lte=30"`
}

type listGuildRAGReferencePoolsResp struct {
	TotalCount int `json:"total_count"`

	RAGReferencePools []*ragReferencePool `json:"rag_reference_pools"`
}

func newlistGuildRAGReferencePoolsRespFromListResult(listResult domain.ListRAGReferencePoolsResult) *listGuildRAGReferencePoolsResp {
	entRAGReferencePools := listResult.GetItems()
	ragReferencePools := make([]*ragReferencePool, len(entRAGReferencePools))
	for i, entRAGReferencePool := range entRAGReferencePools {
		ragReferencePools[i] = newRAGReferencePoolFromEnt(entRAGReferencePool)
	}

	return &listGuildRAGReferencePoolsResp{
		TotalCount:        listResult.GetTotal(),
		RAGReferencePools: ragReferencePools,
	}
}

type ragReferencePool struct {
	ID int `json:"id"`

	Name string `json:"name"`

	Description string `json:"description"`

	CreatedAt time.Time `json:"created_at"`
}

func newRAGReferencePoolFromEnt(entRAGReferencePool *ent.RAGReferencePool) *ragReferencePool {
	return &ragReferencePool{
		ID:          entRAGReferencePool.ID,
		Name:        entRAGReferencePool.Name,
		Description: entRAGReferencePool.Description,
		CreatedAt:   entRAGReferencePool.CreatedAt,
	}
}

type ginAddChatroomRAGReferencePoolReq struct {
	RAGReferencePoolID int `json:"rag_reference_pool_id" binding:"required"`
}

type ginListChatroomRAGReferencePoolsReq struct {
	Page int `form:"page" binding:"required,gte=1"`

	Limit int `form:"limit" binding:"required,gte=1,lte=30"`
}

type listChatroomRAGReferencePoolsResp struct {
	TotalCount int `json:"total_count"`

	RAGReferencePools []*ragReferencePool `json:"rag_reference_pools"`
}

func newlistChatroomRAGReferencePoolsRespFromListResult(listResult domain.ListRAGReferencePoolsResult) *listChatroomRAGReferencePoolsResp {
	entRAGReferencePools := listResult.GetItems()
	ragReferencePools := make([]*ragReferencePool, len(entRAGReferencePools))
	for i, entRAGReferencePool := range entRAGReferencePools {
		ragReferencePools[i] = newRAGReferencePoolFromEnt(entRAGReferencePool)
	}

	return &listChatroomRAGReferencePoolsResp{
		TotalCount:        listResult.GetTotal(),
		RAGReferencePools: ragReferencePools,
	}
}

type ginRemoveChatroomRAGReferencePoolsReq struct {
	RAGReferencePoolIDs []int `json:"rag_reference_pool_ids" binding:"required,min=1"`
}
