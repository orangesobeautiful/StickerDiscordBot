package migrate

import (
	"context"

	"backend/app/config"
	"backend/app/domain"
	meilisearchMigrate "backend/app/server/migrate/meilisearch"

	"github.com/meilisearch/meilisearch-go"
)

type Migrator struct {
	meilisearch meilisearch.ServiceManager

	migrateRepo domain.MigrateRepository

	fullTextSearchConfig config.FullTextSearchDatabase
}

func NewMigrator(
	meilisearchManager meilisearch.ServiceManager,
	migrateRepo domain.MigrateRepository,
	fullTextSearchConfig config.FullTextSearchDatabase,
) *Migrator {
	return &Migrator{
		meilisearch:          meilisearchManager,
		migrateRepo:          migrateRepo,
		fullTextSearchConfig: fullTextSearchConfig,
	}
}

func (m *Migrator) Migrate(ctx context.Context) error {
	meilisearchMigrator := meilisearchMigrate.NewMeilisearchMigrator(
		m.meilisearch,
		m.migrateRepo,
		m.fullTextSearchConfig,
	)

	err := meilisearchMigrator.Migrate(ctx)
	if err != nil {
		return err
	}

	return nil
}
