package domain

import (
	"context"

	"backend/app/ent"
)

type DiscordCommandRepository interface {
	BaseEntRepoInterface

	Add(ctx context.Context, name string, discordID string, sha256Checksum []byte) error
	GetAll(ctx context.Context) ([]*ent.DiscordCommand, error)
	UpdateByName(ctx context.Context, name string, discordID string, sha256Checksum []byte) error
	DeleteByName(ctx context.Context, name string) error
}
