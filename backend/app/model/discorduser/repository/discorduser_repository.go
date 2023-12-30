package repository

import (
	"context"
	"fmt"

	"backend/app/domain"
	"backend/app/ent"
	"backend/app/ent/discorduser"
	"backend/app/pkg/hserr"
)

var _ domain.DiscordUserRepository = (*discordUserRepository)(nil)

type discordUserRepository struct {
	*domain.BaseEntRepo
}

func NewDiscordUser(client *ent.Client) domain.DiscordUserRepository {
	bRepo := domain.NewBaseEntRepo(client)

	return &discordUserRepository{
		BaseEntRepo: bRepo,
	}
}

func (r *discordUserRepository) Create(ctx context.Context, discordID, channelID, name, avatarURL string) (int, error) {
	du, err := r.GetEntClient(ctx).DiscordUser.
		Create().
		SetDiscordID(discordID).
		SetChannelID(channelID).
		SetName(name).
		SetAvatarURL(avatarURL).
		Save(ctx)
	if err != nil {
		return 0, hserr.NewInternalError(err, "create discord user")
	}

	return du.ID, nil
}

func (r *discordUserRepository) FindByDiscordAndChannelID(ctx context.Context, discordID, channelID string) (*ent.DiscordUser, error) {
	du, err := r.GetEntClient(ctx).DiscordUser.
		Query().
		Where(
			discorduser.DiscordID(discordID),
			discorduser.ChannelID(channelID),
		).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, hserr.NewInternalError(err, fmt.Sprintf("query discord user, discordID=%s, channelID=%s", discordID, channelID))
	}

	return du, nil
}
