package domain

import (
	"context"

	"backend/app/ent"
)

type ListChatroomsResult = ListResult[*ent.Chatroom]

type DiscordGuildRepository interface {
	BaseEntRepoInterface
	CreateGuild(ctx context.Context, guildID string) (err error)
	GetGuildByID(ctx context.Context, guildID string) (guild *ent.DiscordGuild, err error)
	CreateGuildChatroom(ctx context.Context, guildID string, name string) (chatroomID int, err error)
	ListGuildChatrooms(ctx context.Context, guildID string, limit, offset int) (result ListChatroomsResult, err error)
	GetGuildActivateChatroomID(ctx context.Context, guildID string) (chatroomID int, err error)
	GetChatroomByID(ctx context.Context, chatroomID int) (chatroom *ent.Chatroom, err error)
	ChangeGuildActivateChatroom(ctx context.Context, guildID string, chatroomID int) (err error)
}

type DiscordGuildUsecase interface {
	RegisterGuild(ctx context.Context, guildID string) (err error)
	CreateGuildChatroom(ctx context.Context, guildID string, name string) (chatroomID int, err error)
	ListGuildChatrooms(ctx context.Context, guildID string, limit, offset int) (result ListChatroomsResult, err error)
	GetGuildActivateChatroom(ctx context.Context, guildID string) (chatroom *ent.Chatroom, err error)
	ChangeGuildActivateChatRoom(ctx context.Context, guildID string, chatroomID int) (err error)
}
