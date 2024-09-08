package meilisearch

import (
	"context"

	"backend/app/config"
	"backend/app/domain"

	"github.com/meilisearch/meilisearch-go"
	"golang.org/x/xerrors"
)

type MeilisearchMigrator struct {
	meilisearch meilisearch.ServiceManager

	migrateRepo domain.MigrateRepository

	fullTextSearchConfig config.FullTextSearchDatabase

	indexNamer *meilisearchIndexName
}

func NewMeilisearchMigrator(
	meilisearchManager meilisearch.ServiceManager,
	migrateRepo domain.MigrateRepository,
	fullTextSearchConfig config.FullTextSearchDatabase,
) *MeilisearchMigrator {
	indexNamePrefix := fullTextSearchConfig.GetMeilisearch().GetIndexPrefix()

	return &MeilisearchMigrator{
		meilisearch:          meilisearchManager,
		migrateRepo:          migrateRepo,
		fullTextSearchConfig: fullTextSearchConfig,
		indexNamer:           newMeilisearchIndexName(indexNamePrefix),
	}
}

type meilisearchIndexName struct {
	prefix string
}

func newMeilisearchIndexName(prefix string) *meilisearchIndexName {
	return &meilisearchIndexName{
		prefix: prefix,
	}
}

func (m *meilisearchIndexName) GetStickerIndex() string {
	return m.prefix + "sticker"
}

func (m *MeilisearchMigrator) Migrate(ctx context.Context) error {
	var err error

	if !m.fullTextSearchConfig.GetToMigrate() {
		return nil
	}

	const targetVersion = 1

	versionInfo, err := m.migrateRepo.GetMeilisearch(ctx)
	if err != nil {
		return xerrors.Errorf("get meilisearch: %w", err)
	}

	currentVersion := versionInfo.Version

	for currentVersion < targetVersion {
		switch currentVersion {
		case 0:
			err = m.migrateMeilisearchTo1(ctx)
		default:
			return xerrors.Errorf("unsupported meilisearch version: %d", currentVersion)
		}

		if err != nil {
			return xerrors.Errorf("migrate meilisearch: %w", err)
		}

		versionInfo, err := m.migrateRepo.GetMeilisearch(ctx)
		if err != nil {
			return xerrors.Errorf("get meilisearch: %w", err)
		}

		currentVersion = versionInfo.Version
	}

	return nil
}

func (m *MeilisearchMigrator) meilisearchToTask(
	ctx context.Context,
	taskFunc func() (*meilisearch.TaskInfo, error),
) error {
	taskInfo, err := taskFunc()
	if err != nil {
		return xerrors.Errorf("task: %w", err)
	}

	task, err := m.meilisearch.WaitForTaskWithContext(ctx, taskInfo.TaskUID, 0)
	if err != nil {
		return xerrors.Errorf("wait for task: %w", err)
	}

	if task.Status != meilisearch.TaskStatusSucceeded {
		return xerrors.Errorf("task failed: %s", task.Error)
	}

	return nil
}

func ptr[T any](v T) *T {
	return &v
}
