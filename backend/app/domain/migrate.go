package domain

import (
	"context"

	"backend/app/ent"
)

type MigrateRepository interface {
	BaseEntRepoInterface
	GetMeilisearch(ctx context.Context) (result *ent.AdditionalMigrate, err error)
	UpdateMeilisearch(ctx context.Context, version int) (err error)
}
