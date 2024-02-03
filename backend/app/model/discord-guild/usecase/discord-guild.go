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
}

func New(
	discordGuildRepo domain.DiscordGuildRepository,
) domain.DiscordGuildUsecase {
	return &discordGuildUsecase{
		discordGuildRepo: discordGuildRepo,
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
