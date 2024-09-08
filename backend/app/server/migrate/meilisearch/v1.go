package meilisearch

import (
	"context"

	"github.com/meilisearch/meilisearch-go"
	"golang.org/x/xerrors"
)

func (m *MeilisearchMigrator) migrateMeilisearchTo1(
	ctx context.Context,
) error {
	err := m.migrateRepo.WithTx(ctx, m.v1Transaction())
	if err != nil {
		return xerrors.Errorf("migrate meilisearch to 1: %w", err)
	}

	return nil
}

func (m *MeilisearchMigrator) v1Transaction() func(context.Context) error {
	return func(ctx context.Context) error {
		var err error

		err = m.migrateRepo.UpdateMeilisearch(ctx, 1)
		if err != nil {
			return xerrors.Errorf("update meilisearch to 1: %w", err)
		}

		err = m.v1StickerIndex(ctx)
		if err != nil {
			return xerrors.Errorf("migrate meilisearch sticker index: %w", err)
		}

		return nil
	}
}

func (m *MeilisearchMigrator) v1StickerIndex(ctx context.Context) error {
	var err error

	stickerHandler := newV1StickerHandler(m, m.indexNamer.Sticker())

	err = stickerHandler.CreateIndex(ctx)
	if err != nil {
		return xerrors.Errorf("create meilisearch sticker index: %w", err)
	}

	err = stickerHandler.UpdateSearchableAttributes(ctx)
	if err != nil {
		return xerrors.Errorf("update meilisearch sticker index: %w", err)
	}

	err = stickerHandler.UpdateSortableAttributes(ctx)
	if err != nil {
		return xerrors.Errorf("update meilisearch sticker index: %w", err)
	}

	err = stickerHandler.UpdateFilterableAttributes(ctx)
	if err != nil {
		return xerrors.Errorf("update meilisearch sticker index: %w", err)
	}

	return nil
}

type v1StickerHandler struct {
	*MeilisearchMigrator

	indexName string

	stickerIndex *meilisearch.IndexResult
}

func newV1StickerHandler(
	migrator *MeilisearchMigrator,
	indexName string,
) *v1StickerHandler {
	return &v1StickerHandler{
		MeilisearchMigrator: migrator,
		indexName:           indexName,
	}
}

func (h *v1StickerHandler) CreateIndex(ctx context.Context) error {
	var err error

	err = h.meilisearchToTask(ctx,
		func() (*meilisearch.TaskInfo, error) {
			return h.meilisearch.CreateIndexWithContext(
				ctx,
				&meilisearch.IndexConfig{
					Uid:        h.indexNamer.Sticker(),
					PrimaryKey: "id",
				},
			)
		},
	)
	if err != nil {
		return xerrors.Errorf("create meilisearch sticker index: %w", err)
	}

	h.stickerIndex, err = h.meilisearch.GetIndexWithContext(ctx, h.indexNamer.Sticker())
	if err != nil {
		return xerrors.Errorf("get meilisearch sticker index: %w", err)
	}

	return nil
}

func (h *v1StickerHandler) UpdateSearchableAttributes(ctx context.Context) error {
	err := h.meilisearchToTask(ctx, func() (*meilisearch.TaskInfo, error) {
		return h.stickerIndex.UpdateSearchableAttributesWithContext(
			ctx,
			ptr([]string{"name"}),
		)
	})
	if err != nil {
		return xerrors.Errorf("update meilisearch sticker index: %w", err)
	}

	return nil
}

func (h *v1StickerHandler) UpdateSortableAttributes(ctx context.Context) error {
	err := h.meilisearchToTask(ctx, func() (*meilisearch.TaskInfo, error) {
		return h.stickerIndex.UpdateSortableAttributesWithContext(
			ctx,
			ptr([]string{"created_at"}),
		)
	})
	if err != nil {
		return xerrors.Errorf("update meilisearch sticker index: %w", err)
	}

	return nil
}

func (h *v1StickerHandler) UpdateFilterableAttributes(ctx context.Context) error {
	err := h.meilisearchToTask(ctx, func() (*meilisearch.TaskInfo, error) {
		return h.stickerIndex.UpdateFilterableAttributesWithContext(
			ctx,
			ptr([]string{"discord_guild_sticker", "created_at"}),
		)
	})
	if err != nil {
		return xerrors.Errorf("update meilisearch sticker index: %w", err)
	}

	return nil
}
