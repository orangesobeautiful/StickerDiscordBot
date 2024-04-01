package usecase

import (
	"context"

	"backend/app/domain"

	"golang.org/x/xerrors"
)

func (u *discordGuildUsecase) CreateRAGReferencePool(
	ctx context.Context, guildID, name, description string,
) (id int, err error) {
	id, err = u.chatRepo.CreateRAGReferencePool(ctx, guildID, name, description)
	if err != nil {
		return 0, xerrors.Errorf("create rag reference pool: %w", err)
	}

	return id, nil
}

func (u *discordGuildUsecase) ListRAGReferencePools(
	ctx context.Context, guildID string, limit, offset int,
) (result domain.ListRAGReferencePoolsResult, err error) {
	result, err = u.chatRepo.ListRAGReferencePools(ctx, guildID, limit, offset)
	if err != nil {
		return domain.ListRAGReferencePoolsResult{}, xerrors.Errorf("list rag reference pools: %w", err)
	}

	return result, nil
}

func (u *discordGuildUsecase) IsGuildOwnRAGReferencePool(
	ctx context.Context, guildID string, ragReferencePoolID int,
) (isOwn bool, err error) {
	ragReferencePool, err := u.chatRepo.GetRAGReferencePoolWithGuildByID(ctx, ragReferencePoolID)
	if err != nil {
		return false, xerrors.Errorf("get chatroom by id: %w", err)
	}

	ragReferencePoolGuild, err := ragReferencePool.Edges.GuildOrErr()
	if err != nil {
		return false, xerrors.Errorf("get chatroom guild: %w", err)
	}

	if ragReferencePoolGuild.ID != guildID {
		return false, nil
	}

	return true, nil
}

func (u *discordGuildUsecase) IsGuildOwnRAGReferenceText(
	ctx context.Context, guildID string, ragReferenceTextID int,
) (isOwn bool, err error) {
	ragReferenceText, err := u.chatRepo.GetRAGReferenceTextWithGuildByID(ctx, ragReferenceTextID)
	if err != nil {
		return false, xerrors.Errorf("get rag reference text by id: %w", err)
	}

	ragReferencePool, err := ragReferenceText.Edges.RefOrErr()
	if err != nil {
		return false, xerrors.Errorf("get rag reference text guild: %w", err)
	}
	guild, err := ragReferencePool.Edges.GuildOrErr()
	if err != nil {
		return false, xerrors.Errorf("get rag reference text guild: %w", err)
	}

	if guild.ID != guildID {
		return false, nil
	}

	return true, nil
}
