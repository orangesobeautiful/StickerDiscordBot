package domain

import (
	"context"

	"backend/app/ent"
)

type DiscordUserRepository interface {
	BaseEntRepoInterface
	Create(ctx context.Context, discordID, channelID, name, avatarURL string) (int, error)
	FindByDiscordAndChannelID(ctx context.Context, discordID, channelID string) (*ent.DiscordUser, error)
}

type DiscordUserUsecase interface {
	Create(ctx context.Context, discordID, channelID, name, avatarURL string) (int, error)
	FindByDiscordAndChannelID(ctx context.Context, discordID, channelID string) (*ent.DiscordUser, error)
}

type DiscordWebLoginVerification struct {
	code          string
	userDiscordID string
	userChannelID string
}

func NewDiscordWebLoginLoginVerification(code, userDiscordID, userChannelID string) *DiscordWebLoginVerification {
	return &DiscordWebLoginVerification{
		code:          code,
		userDiscordID: userDiscordID,
		userChannelID: userChannelID,
	}
}

func (l *DiscordWebLoginVerification) GetCode() string {
	return l.code
}

func (l *DiscordWebLoginVerification) GetUserDiscordID() string {
	return l.userDiscordID
}

type DiscordWebLoginVerificationRepository interface {
	Create(ctx context.Context, verifyCode, dcUserID, dcChannelID string) error
	FindByCode(ctx context.Context, code string) (*DiscordWebLoginVerification, error)
}

type DiscordWebLoginVerificationUsecase interface {
	Create(ctx context.Context, verifyCode, dcUserID, dcChannelID string) error
	FindByCode(ctx context.Context, code string) (*DiscordWebLoginVerification, error)
}
