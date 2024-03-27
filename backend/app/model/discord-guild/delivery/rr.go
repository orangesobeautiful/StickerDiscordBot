package delivery

import (
	"fmt"
	"time"

	"backend/app/domain"
	"backend/app/ent"
	discordcommand "backend/app/pkg/discord-command"

	"github.com/bwmarrin/discordgo"
)

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
