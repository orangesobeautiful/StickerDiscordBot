package repository

import (
	"context"

	"backend/app/domain"
	"backend/app/ent"
	"backend/app/ent/chatroom"
	"backend/app/ent/discordguild"
	"backend/app/pkg/hserr"
)

var _ domain.DiscordGuildRepository = (*discordGuildRepository)(nil)

type discordGuildRepository struct {
	*domain.BaseEntRepo
}

func New(client *ent.Client) domain.DiscordGuildRepository {
	bRepo := domain.NewBaseEntRepo(client)

	return &discordGuildRepository{
		BaseEntRepo: bRepo,
	}
}

func (r *discordGuildRepository) CreateGuild(ctx context.Context, guildID string) (err error) {
	_, err = r.GetEntClient(ctx).DiscordGuild.
		Create().
		SetID(guildID).
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			err = domain.NewAlreadyExistsError("discord guild")
			return hserr.NewInternalError(err, "create discord guild")
		}

		return hserr.NewInternalError(err, "create discord guild")
	}

	return nil
}

func (r *discordGuildRepository) GetGuildByID(ctx context.Context, guildID string) (guild *ent.DiscordGuild, err error) {
	guild, err = r.GetEntClient(ctx).DiscordGuild.
		Query().
		Where(
			discordguild.ID(guildID),
		).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NewNotFoundError("discord guild")
			return nil, hserr.NewInternalError(err, "get discord guild by id")
		}
		return nil, hserr.NewInternalError(err, "get discord guild by id")
	}

	return guild, nil
}

func (r *discordGuildRepository) CreateGuildChatroom(
	ctx context.Context, guildID string, name string,
) (chatroomID int, err error) {
	newChatroom, err := r.GetEntClient(ctx).Chatroom.
		Create().
		SetGuildID(guildID).
		SetName(name).
		Save(ctx)
	if err != nil {
		return 0, hserr.NewInternalError(err, "create chatroom")
	}

	return newChatroom.ID, nil
}

func (r *discordGuildRepository) GetChatroomByID(ctx context.Context, chatroomID int) (chatroomResult *ent.Chatroom, err error) {
	chatroomResult, err = r.GetEntClient(ctx).Chatroom.
		Query().
		Where(
			chatroom.ID(chatroomID),
		).
		Only(ctx)
	if err != nil {
		return nil, hserr.NewInternalError(err, "get chatroom by id")
	}

	return chatroomResult, nil
}

func (r *discordGuildRepository) ListGuildChatrooms(
	ctx context.Context, guildID string, limit, offset int,
) (result domain.ListChatroomsResult, err error) {
	queryFilter := r.GetEntClient(ctx).Chatroom.
		Query().
		Where(
			chatroom.HasGuildWith(discordguild.ID(guildID)),
		)

	total, err := queryFilter.Clone().Count(ctx)
	if err != nil {
		return result, hserr.NewInternalError(err, "query chatroom count")
	}

	chatrooms, err := queryFilter.
		Limit(limit).
		Offset(offset).
		All(ctx)
	if err != nil {
		return domain.ListChatroomsResult{}, hserr.NewInternalError(err, "list chatrooms")
	}

	result = domain.NewListResult(total, chatrooms)
	return result, nil
}

func (r *discordGuildRepository) GetGuildActivateChatroomID(ctx context.Context, guildID string) (chatroomID int, err error) {
	discordGuild, err := r.GetEntClient(ctx).DiscordGuild.
		Query().
		Where(
			discordguild.ID(guildID),
		).
		WithActivateChatroom().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			err = domain.NewNotFoundError("discord guild")
			return 0, hserr.NewInternalError(err, "get discord guild by id")
		}

		return 0, hserr.NewInternalError(err, "get activate chatroom")
	}

	chatroomID = discordGuild.Edges.ActivateChatroom.ID
	return chatroomID, nil
}

func (r *discordGuildRepository) ChangeGuildActivateChatroom(ctx context.Context, guildID string, chatroomID int) (err error) {
	_, err = r.GetEntClient(ctx).DiscordGuild.
		UpdateOneID(guildID).
		SetActivateChatroomID(chatroomID).
		Save(ctx)
	if err != nil {
		return hserr.NewInternalError(err, "change activated chatroom")
	}

	return nil
}
