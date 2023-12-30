package usecase

import (
	"context"

	"backend/app/domain"

	"golang.org/x/xerrors"
)

var _ domain.DiscordWebLoginVerificationUsecase = (*dcWebLoginVerifyUsecase)(nil)

type (
	dcWebLoginVerifyUsecase            = discordWebLoginVerificationUsecase
	discordWebLoginVerificationUsecase struct {
		dcWebLoginVerifyRepo domain.DiscordWebLoginVerificationRepository
	}
)

func NewDCWebUsecase(
	dcWebLoginVerifyRepo domain.DiscordWebLoginVerificationRepository,
) domain.DiscordWebLoginVerificationUsecase {
	return &dcWebLoginVerifyUsecase{
		dcWebLoginVerifyRepo: dcWebLoginVerifyRepo,
	}
}

func (s *dcWebLoginVerifyUsecase) Create(ctx context.Context, verifyCode, dcUserID, dcChannelID string) (err error) {
	err = s.dcWebLoginVerifyRepo.Create(ctx, verifyCode, dcUserID, dcChannelID)
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	return nil
}

func (s *dcWebLoginVerifyUsecase) FindByCode(ctx context.Context, code string) (result *domain.DiscordWebLoginVerification, err error) {
	result, err = s.dcWebLoginVerifyRepo.FindByCode(ctx, code)
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return result, nil
}
