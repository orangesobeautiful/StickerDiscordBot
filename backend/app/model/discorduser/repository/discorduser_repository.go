package repository

import (
	"context"
	"fmt"

	"backend/app/domain"
	"backend/app/ent"
	"backend/app/ent/discorduser"
	"backend/app/pkg/hserr"

	"entgo.io/ent/dialect/sql"
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

func (r *discordUserRepository) Upsert(ctx context.Context, discordID, guildID, name, avatarURL string) (id int, err error) {
	id, err = r.GetEntClient(ctx).DiscordUser.
		Create().
		SetDiscordID(discordID).
		SetGuildID(guildID).
		SetName(name).
		SetAvatarURL(avatarURL).
		OnConflict(
			sql.ConflictColumns(discorduser.FieldDiscordID, discorduser.FieldGuildID),
			sql.ResolveWithNewValues(),
		).
		Update(func(duu *ent.DiscordUserUpsert) {
			duu.SetName(name)
			duu.SetAvatarURL(avatarURL)
		}).
		ID(ctx)
	if err != nil {
		return 0, hserr.NewInternalError(err, "upsert discord user")
	}

	return id, nil
}

func (r *discordUserRepository) FindByDiscordAndGuildlID(ctx context.Context, discordID, guildID string) (*ent.DiscordUser, error) {
	du, err := r.GetEntClient(ctx).DiscordUser.
		Query().
		Where(
			discorduser.DiscordID(discordID),
			discorduser.GuildID(guildID),
		).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, hserr.NewInternalError(err, fmt.Sprintf("query discord user, discordID=%s, guildID=%s", discordID, guildID))
	}

	return du, nil
}
