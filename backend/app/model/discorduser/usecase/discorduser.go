package usecase

import (
	"context"

	"backend/app/domain"
	"backend/app/ent"

	"golang.org/x/xerrors"
)

var _ domain.DiscordUserUsecase = (*discordUserUsecase)(nil)

type discordUserUsecase struct {
	discordUserRepository domain.DiscordUserRepository
}

func NewDiscordUser(discordUserRepo domain.DiscordUserRepository) domain.DiscordUserUsecase {
	return &discordUserUsecase{
		discordUserRepository: discordUserRepo,
	}
}

func (s *discordUserUsecase) Create(ctx context.Context, discordID, channelID, name, avatarURL string) (id int, err error) {
	id, err = s.discordUserRepository.Create(ctx, discordID, channelID, name, avatarURL)
	if err != nil {
		return 0, xerrors.Errorf(": %w", err)
	}

	return id, nil
}

func (s *discordUserUsecase) FindByDiscordAndChannelID(
	ctx context.Context, discordID, channelID string,
) (user *ent.DiscordUser, err error) {
	user, err = s.discordUserRepository.FindByDiscordAndChannelID(ctx, discordID, channelID)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return user, nil
}
