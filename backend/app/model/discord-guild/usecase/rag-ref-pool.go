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
