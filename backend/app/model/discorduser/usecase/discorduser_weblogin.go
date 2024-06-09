package usecase

import (
	"context"
	"net/http"

	"backend/app/domain"
	"backend/app/ent"
	"backend/app/pkg/hserr"
	"backend/app/utils"

	"github.com/google/uuid"
	"golang.org/x/xerrors"
)

var _ domain.DiscordWebLoginVerificationUsecase = (*dcWebLoginVerifyUsecase)(nil)

type (
	dcWebLoginVerifyUsecase            = discordWebLoginVerificationUsecase
	discordWebLoginVerificationUsecase struct {
		dcWebLoginVerifyRepo  domain.DiscordWebLoginVerificationRepository
		discordUserRepository domain.DiscordUserRepository
	}
)

func NewDCWebUsecase(
	dcWebLoginVerifyRepo domain.DiscordWebLoginVerificationRepository,
	discordUserRepo domain.DiscordUserRepository,
) domain.DiscordWebLoginVerificationUsecase {
	return &dcWebLoginVerifyUsecase{
		dcWebLoginVerifyRepo:  dcWebLoginVerifyRepo,
		discordUserRepository: discordUserRepo,
	}
}

var verifyCodeStringSet = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func (s *dcWebLoginVerifyUsecase) CreateRandomLoginCode(ctx context.Context) (code string, err error) {
	const verifyCodeLength = 6
	code = utils.RandString(verifyCodeStringSet, verifyCodeLength)

	err = s.dcWebLoginVerifyRepo.CreateLoginCode(ctx, code)
	if err != nil {
		return "", xerrors.Errorf("create code: %w", err)
	}

	return code, nil
}

func (s *dcWebLoginVerifyUsecase) VerifyLoginCode(
	ctx context.Context, code, userDiscordID, userGuildlID, name, avatarURL string,
) (err error) {
	err = s.verifyLoginCodeCheck(ctx, code)
	if err != nil {
		return xerrors.Errorf("verify login code check: %w", err)
	}

	err = s.dcWebLoginVerifyRepo.UpdateDiscordUserInfoByCode(
		ctx, code, userDiscordID, userGuildlID, name, avatarURL,
	)
	if err != nil {
		return xerrors.Errorf("update dc user info: %w", err)
	}

	return nil
}

func (s *dcWebLoginVerifyUsecase) verifyLoginCodeCheck(ctx context.Context, code string) (err error) {
	verifyInfo, err := s.dcWebLoginVerifyRepo.FindLoginCodeByCode(ctx, code)
	if err != nil {
		return xerrors.Errorf("find by code: %w", err)
	}
	if verifyInfo == nil {
		return hserr.New(http.StatusBadRequest, "verify code not found")
	}

	if verifyInfo.IsVerified() {
		return hserr.New(http.StatusBadRequest, "already verified")
	}

	return nil
}

func (s *dcWebLoginVerifyUsecase) CheckLoginCode(ctx context.Context, code string) (newSessionID uuid.UUID, loggedIn bool, err error) {
	verifyResult, err := s.dcWebLoginVerifyRepo.FindLoginCodeByCode(ctx, code)
	if err != nil {
		return uuid.Nil, false, xerrors.Errorf("find by code: %w", err)
	}
	if verifyResult == nil {
		return uuid.Nil, false, hserr.New(http.StatusBadRequest, "verify code is not valid")
	}

	if !verifyResult.IsVerified() {
		return uuid.Nil, false, nil
	}

	newSessionID, err = s.createNewLoginSessionByVerifyResult(ctx, verifyResult)
	if err != nil {
		return uuid.Nil, false, xerrors.Errorf("create login session: %w", err)
	}

	_ = s.dcWebLoginVerifyRepo.DeleteLoginCode(ctx, code)
	return newSessionID, true, nil
}

func (s *dcWebLoginVerifyUsecase) createNewLoginSessionByVerifyResult(
	ctx context.Context, verifyResult *domain.DiscordWebLoginVerification,
) (newSessionID uuid.UUID, err error) {
	err = s.dcWebLoginVerifyRepo.WithTx(ctx, func(ctx context.Context) error {
		var userID int
		userID, err = s.upsertDCUserInfoByLoginVerifyResult(ctx, verifyResult)
		if err != nil {
			return xerrors.Errorf("upsert dc user info: %w", err)
		}

		newSessionID, err = s.dcWebLoginVerifyRepo.CreateLoginSession(ctx, userID)
		if err != nil {
			return xerrors.Errorf("create login session: %w", err)
		}

		err = s.dcWebLoginVerifyRepo.DeleteLoginCode(ctx, verifyResult.GetCode())
		if err != nil {
			return xerrors.Errorf("delete login code: %w", err)
		}

		return nil
	})
	if err != nil {
		return uuid.Nil, xerrors.Errorf(": %w", err)
	}

	return newSessionID, nil
}

func (s *dcWebLoginVerifyUsecase) upsertDCUserInfoByLoginVerifyResult(
	ctx context.Context, verifyResult *domain.DiscordWebLoginVerification,
) (dcUserID int, err error) {
	return s.discordUserRepository.Upsert(ctx,
		verifyResult.GetUserDiscordID(),
		verifyResult.GetUserGuildID(),
		verifyResult.GetUserName(),
		verifyResult.GetUserAvatarURL(),
	)
}

func (s dcWebLoginVerifyUsecase) GetDiscordUserBySessionID(ctx context.Context, sessionID uuid.UUID) (dcUser *ent.DiscordUser, err error) {
	dcUser, err = s.dcWebLoginVerifyRepo.GetDiscordUserBySessionID(ctx, sessionID)
	if err != nil {
		return nil, xerrors.Errorf("get dc user by session: %w", err)
	}

	return dcUser, nil
}
