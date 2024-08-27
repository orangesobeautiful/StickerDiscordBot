package repository

import (
	"context"

	"backend/app/domain"
	"backend/app/ent"
	"backend/app/ent/discordcommand"
	"backend/app/pkg/hserr"
)

var _ domain.DiscordCommandRepository = (*discordCommandRepository)(nil)

type discordCommandRepository struct {
	*domain.BaseEntRepo
}

func New(client *ent.Client) domain.DiscordCommandRepository {
	bRepo := domain.NewBaseEntRepo(client)

	return &discordCommandRepository{
		BaseEntRepo: bRepo,
	}
}

func (r *discordCommandRepository) Add(ctx context.Context, name string, discordID string, sha256Checksum []byte) (err error) {
	_, err = r.GetEntClient(ctx).DiscordCommand.
		Create().
		SetName(name).
		SetDiscordID(discordID).
		SetSha256Checksum(sha256Checksum).
		Save(ctx)
	if err != nil {
		return hserr.NewInternalError(err, "create discord command")
	}

	return nil
}

func (r *discordCommandRepository) GetAll(ctx context.Context) (commands []*ent.DiscordCommand, err error) {
	commands, err = r.GetEntClient(ctx).DiscordCommand.
		Query().
		All(ctx)
	if err != nil {
		return nil, hserr.NewInternalError(err, "query discord commands")
	}

	return commands, nil
}

func (r *discordCommandRepository) UpdateByName(ctx context.Context, name string, discordID string, sha256Checksum []byte) (err error) {
	_, err = r.GetEntClient(ctx).DiscordCommand.
		Update().
		Where(
			discordcommand.Name(name),
		).
		SetDiscordID(discordID).
		SetSha256Checksum(sha256Checksum).
		Save(ctx)
	if err != nil {
		return hserr.NewInternalError(err, "update discord command")
	}

	return nil
}

func (r *discordCommandRepository) DeleteByName(ctx context.Context, name string) (err error) {
	_, err = r.GetEntClient(ctx).DiscordCommand.
		Delete().
		Where(
			discordcommand.Name(name),
		).
		Exec(ctx)
	if err != nil {
		return hserr.NewInternalError(err, "delete discord command")
	}

	return nil
}
