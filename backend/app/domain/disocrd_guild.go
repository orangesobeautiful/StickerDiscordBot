package domain

import (
	"context"

	"backend/app/ent"
)

type (
	ListChatroomsResult         = ListResult[*ent.Chatroom]
	ListRAGReferencePoolsResult = ListResult[*ent.RAGReferencePool]
	ListRAGReferenceTextsResult = ListResult[*ent.RAGReferenceText]
)

type DiscordGuildRepository interface {
	BaseEntRepoInterface
	CreateGuild(ctx context.Context, guildID string) (err error)
	GetGuildByID(ctx context.Context, guildID string) (guild *ent.DiscordGuild, err error)

	CreateGuildChatroom(ctx context.Context, guildID string, name string) (chatroomID int, err error)
	ListGuildChatrooms(ctx context.Context, guildID string, limit, offset int) (result ListChatroomsResult, err error)
	GetChatroomByID(ctx context.Context, chatroomID int) (chatroom *ent.Chatroom, err error)
	GetChatroomWithGuildByID(ctx context.Context, chatroomID int) (chatroom *ent.Chatroom, err error)
	RemoveGuildChatroom(ctx context.Context, chatroomID int) (err error)

	GetGuildActivateChatroomID(ctx context.Context, guildID string) (chatroomID int, err error)
	IsChatroomActivate(ctx context.Context, chatroomID int) (isActivate bool, err error)
	ChangeGuildActivateChatroom(ctx context.Context, guildID string, chatroomID int) (err error)

	AddChatroomRAGReferencePool(ctx context.Context, chatroomID int, ragReferencePoolID int) (err error)
	GetAllChatroomRAGReferencePools(ctx context.Context, chatroomID int) (result []*ent.RAGReferencePool, err error)
	ListChatroomRAGReferencePools(ctx context.Context, chatroomID int, limit, offset int) (result ListRAGReferencePoolsResult, err error)
	RemoveChatroomRAGReferencePools(ctx context.Context, chatroomID int, ragReferencePoolID []int) (err error)
}

type DiscordGuildUsecase interface {
	RegisterGuild(ctx context.Context, guildID string) (err error)

	CreateGuildChatroom(ctx context.Context, guildID string, name string) (chatroomID int, err error)
	ListGuildChatrooms(ctx context.Context, guildID string, limit, offset int) (result ListChatroomsResult, err error)
	IsGuildOwnChatroom(ctx context.Context, guildID string, chatroomID int) (isOwn bool, err error)
	RemoveGuildChatroom(ctx context.Context, chatroomID int) (err error)

	GetGuildActivateChatroom(ctx context.Context, guildID string) (chatroom *ent.Chatroom, err error)
	ChangeGuildActivateChatRoom(ctx context.Context, guildID string, chatroomID int) (err error)

	CreateRAGReferencePool(ctx context.Context, guildID, name, description string) (id int, err error)
	ListRAGReferencePools(ctx context.Context, guildID string, limit, offset int) (result ListRAGReferencePoolsResult, err error)

	AddChatroomRAGReferencePool(ctx context.Context, chatroomID int, ragReferencePoolID int) (err error)
	GetAllChatroomRAGReferencePools(ctx context.Context, chatroomID int) (result []*ent.RAGReferencePool, err error)
	ListChatroomRAGReferencePools(ctx context.Context, chatroomID int, limit, offset int) (result ListRAGReferencePoolsResult, err error)
	RemoveChatroomRAGReferencePools(ctx context.Context, chatroomID int, ragReferencePoolIDs []int) (err error)
}
