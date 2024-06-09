package domain

import (
	"context"

	"backend/app/ent"

	"github.com/google/uuid"
)

type DiscordUserRepository interface {
	BaseEntRepoInterface
	Upsert(ctx context.Context, userDiscordID, userGuildlID, name, avatarURL string) (id int, err error)
}

type DiscordWebLoginVerification struct {
	code          string
	userDiscordID string
	userGuildID   string
	userName      string
	userAvatarURL string
}

func NewDiscordWebLoginLoginVerification(
	code string,
	userDiscordID string, userGuildID string, userName string, userAvatarURL string,
) *DiscordWebLoginVerification {
	return &DiscordWebLoginVerification{
		code:          code,
		userDiscordID: userDiscordID,
		userGuildID:   userGuildID,
		userName:      userName,
		userAvatarURL: userAvatarURL,
	}
}

func (l *DiscordWebLoginVerification) GetCode() string {
	return l.code
}

func (l *DiscordWebLoginVerification) GetUserDiscordID() string {
	return l.userDiscordID
}

func (l *DiscordWebLoginVerification) GetUserGuildID() string {
	return l.userGuildID
}

func (l *DiscordWebLoginVerification) GetUserName() string {
	return l.userName
}

func (l *DiscordWebLoginVerification) GetUserAvatarURL() string {
	return l.userAvatarURL
}

func (l *DiscordWebLoginVerification) IsVerified() bool {
	return l.userDiscordID != "" && l.userGuildID != ""
}

type DiscordWebLoginVerificationRepository interface {
	BaseEntRepoInterface
	CreateLoginCode(ctx context.Context, verifyCode string) error
	FindLoginCodeByCode(ctx context.Context, verifyCode string) (*DiscordWebLoginVerification, error)
	UpdateDiscordUserInfoByCode(ctx context.Context, verifyCode, userDiscordID, userGuildlID, name, avatarURL string) error
	DeleteLoginCode(ctx context.Context, verifyCode string) error

	CreateLoginSession(ctx context.Context, dcUserID int) (sessionID uuid.UUID, err error)
	GetDiscordUserBySessionID(ctx context.Context, sessionID uuid.UUID) (*ent.DiscordUser, error)
}

type DiscordWebLoginVerificationUsecase interface {
	CreateRandomLoginCode(ctx context.Context) (code string, err error)
	VerifyLoginCode(ctx context.Context, code, userDiscordID, userGuildlID, name, avatarURL string) (err error)
	CheckLoginCode(ctx context.Context, code string) (newSessionID uuid.UUID, loggedIn bool, err error)

	GetDiscordUserBySessionID(ctx context.Context, sessionID uuid.UUID) (dcUser *ent.DiscordUser, err error)
}
