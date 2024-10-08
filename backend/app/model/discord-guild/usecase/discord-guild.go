package usecase

import (
	"context"
	"net/http"

	"backend/app/domain"
	"backend/app/ent"
	"backend/app/pkg/hserr"

	"golang.org/x/xerrors"
)

var _ domain.DiscordGuildUsecase = (*discordGuildUsecase)(nil)

type discordGuildUsecase struct {
	discordGuildRepo domain.DiscordGuildRepository

	stickerRepo domain.StickerRepository

	chatRepo domain.ChatRepository
}

func New(
	discordGuildRepo domain.DiscordGuildRepository,
	stickerRepo domain.StickerRepository,
	chatRepo domain.ChatRepository,
) domain.DiscordGuildUsecase {
	return &discordGuildUsecase{
		discordGuildRepo: discordGuildRepo,

		stickerRepo: stickerRepo,

		chatRepo: chatRepo,
	}
}

func (u *discordGuildUsecase) RegisterGuild(ctx context.Context, guildID string) (err error) {
	err = u.discordGuildRepo.CreateGuild(ctx, guildID)
	if err != nil {
		if domain.IsAlreadyExistsError(err) {
			return hserr.New(http.StatusConflict, "guild already exists")
		}

		return xerrors.Errorf("create guild: %w", err)
	}

	chatroomID, err := u.discordGuildRepo.CreateGuildChatroom(ctx, guildID, "default")
	if err != nil {
		return xerrors.Errorf("create guild chatroom: %w", err)
	}

	err = u.discordGuildRepo.ChangeGuildActivateChatroom(ctx, guildID, chatroomID)
	if err != nil {
		return xerrors.Errorf("change activate chatroom: %w", err)
	}

	return nil
}

func (u *discordGuildUsecase) IsGuildOwnSticker(
	ctx context.Context, guildID string, stickerID int,
) (isOwn bool, err error) {
	sticker, err := u.stickerRepo.GetStickerWithGuildByID(ctx, stickerID)
	if err != nil {
		return false, xerrors.Errorf("get sticker by id: %w", err)
	}

	guild, err := sticker.Edges.GuildOrErr()
	if err != nil {
		return false, hserr.NewInternalError(err, "get sticker guild")
	}

	if guild.ID != guildID {
		return false, nil
	}

	return true, nil
}

func (u *discordGuildUsecase) IsGuildOwnStickerImage(
	ctx context.Context, guildID string, stickerImageID int,
) (isOwn bool, err error) {
	stickerImage, err := u.stickerRepo.GetStickerImageWithGuildByID(ctx, stickerImageID)
	if err != nil {
		return false, xerrors.Errorf("get sticker image by id: %w", err)
	}

	sticker, err := stickerImage.Edges.StickerOrErr()
	if err != nil {
		return false, hserr.NewInternalError(err, "get sticker image sticker")
	}

	guild, err := sticker.Edges.GuildOrErr()
	if err != nil {
		return false, hserr.NewInternalError(err, "get sticker image guild")
	}

	if guild.ID != guildID {
		return false, nil
	}

	return true, nil
}

func (u *discordGuildUsecase) CreateGuildChatroom(
	ctx context.Context, guildID string, name string,
) (chatroomID int, err error) {
	chatroomID, err = u.discordGuildRepo.CreateGuildChatroom(ctx, guildID, name)
	if err != nil {
		return 0, xerrors.Errorf("create chatroom: %w", err)
	}

	return chatroomID, nil
}

func (u *discordGuildUsecase) ListGuildChatrooms(
	ctx context.Context, guildID string, limit, offset int,
) (result domain.ListChatroomsResult, err error) {
	result, err = u.discordGuildRepo.ListGuildChatrooms(ctx, guildID, limit, offset)
	if err != nil {
		return domain.ListChatroomsResult{}, xerrors.Errorf("list chatrooms: %w", err)
	}

	return result, nil
}

func (u *discordGuildUsecase) IsGuildOwnChatroom(
	ctx context.Context, guildID string, chatroomID int,
) (isOwn bool, err error) {
	chatroom, err := u.discordGuildRepo.GetChatroomWithGuildByID(ctx, chatroomID)
	if err != nil {
		return false, xerrors.Errorf("get chatroom by id: %w", err)
	}

	chatroomGuild, err := chatroom.Edges.GuildOrErr()
	if err != nil {
		return false, xerrors.Errorf("get chatroom guild: %w", err)
	}

	if chatroomGuild.ID != guildID {
		return false, nil
	}

	return true, nil
}

func (u *discordGuildUsecase) RemoveGuildChatroom(
	ctx context.Context, chatroomID int,
) (err error) {
	isActivate, err := u.discordGuildRepo.IsChatroomActivate(ctx, chatroomID)
	if err != nil {
		return xerrors.Errorf("is chatroom activate: %w", err)
	}
	if isActivate {
		return hserr.New(http.StatusForbidden, "cannot remove activate chatroom")
	}

	err = u.discordGuildRepo.RemoveGuildChatroom(ctx, chatroomID)
	if err != nil {
		return xerrors.Errorf("remove chatroom: %w", err)
	}

	return nil
}

func (u *discordGuildUsecase) GetGuildActivateChatroom(ctx context.Context, guildID string) (chatroom *ent.Chatroom, err error) {
	var chatroomID int
	chatroomID, err = u.discordGuildRepo.GetGuildActivateChatroomID(ctx, guildID)
	if err != nil {
		if !domain.IsNotFoundError(err) {
			return nil, xerrors.Errorf("get activate chatroom id: %w", err)
		}

		chatroomID, err = u.initGuildAndAcvitateChatroom(ctx, guildID)
		if err != nil {
			return nil, xerrors.Errorf("init guild and activate chatroom: %w", err)
		}
	}

	chatroom, err = u.discordGuildRepo.GetChatroomByID(ctx, chatroomID)
	if err != nil {
		return nil, xerrors.Errorf("get chatroom by id: %w", err)
	}

	return chatroom, nil
}

func (u *discordGuildUsecase) initGuildAndAcvitateChatroom(ctx context.Context, guildID string) (chatroomID int, err error) {
	err = u.discordGuildRepo.CreateGuild(ctx, guildID)
	if err != nil {
		return 0, xerrors.Errorf("create guild: %w", err)
	}

	chatroomID, err = u.discordGuildRepo.CreateGuildChatroom(ctx, guildID, "default")
	if err != nil {
		return 0, xerrors.Errorf("create guild chatroom: %w", err)
	}

	err = u.ChangeGuildActivateChatRoom(ctx, guildID, chatroomID)
	if err != nil {
		return 0, xerrors.Errorf("change activate chatroom: %w", err)
	}

	return chatroomID, nil
}

func (u *discordGuildUsecase) ChangeGuildActivateChatRoom(
	ctx context.Context, guildID string, chatroomID int,
) (err error) {
	err = u.discordGuildRepo.ChangeGuildActivateChatroom(ctx, guildID, chatroomID)
	if err != nil {
		return xerrors.Errorf("change activate chatroom: %w", err)
	}

	return nil
}

func (u *discordGuildUsecase) AddChatroomRAGReferencePool(
	ctx context.Context, chatroomID int, ragReferencePoolID int,
) (err error) {
	err = u.discordGuildRepo.AddChatroomRAGReferencePool(ctx, chatroomID, ragReferencePoolID)
	if err != nil {
		return xerrors.Errorf("add chatroom rag reference pool: %w", err)
	}

	return nil
}

func (u *discordGuildUsecase) GetAllChatroomRAGReferencePools(
	ctx context.Context, chatroomID int,
) (result []*ent.RAGReferencePool, err error) {
	result, err = u.discordGuildRepo.GetAllChatroomRAGReferencePools(ctx, chatroomID)
	if err != nil {
		return nil, xerrors.Errorf("get all chatroom rag reference pools: %w", err)
	}

	return result, nil
}

func (u *discordGuildUsecase) ListChatroomRAGReferencePools(
	ctx context.Context, chatroomID int, limit, offset int,
) (result domain.ListRAGReferencePoolsResult, err error) {
	result, err = u.discordGuildRepo.ListChatroomRAGReferencePools(ctx, chatroomID, limit, offset)
	if err != nil {
		return domain.ListRAGReferencePoolsResult{}, xerrors.Errorf("list chatroom rag reference pools: %w", err)
	}

	return result, nil
}

func (u *discordGuildUsecase) RemoveChatroomRAGReferencePools(
	ctx context.Context, chatroomID int, ragReferencePoolIDs []int,
) (err error) {
	err = u.discordGuildRepo.RemoveChatroomRAGReferencePools(ctx, chatroomID, ragReferencePoolIDs)
	if err != nil {
		return xerrors.Errorf("remove chatroom rag reference pool: %w", err)
	}

	return nil
}
