package repository

import (
	"context"

	"backend/app/domain"
	"backend/app/ent"
	"backend/app/ent/additionalmigrate"
	"backend/app/pkg/hserr"

	"entgo.io/ent/dialect/sql"
)

const (
	migrateNameMeilisearch string = "meilisearch"
)

var _ domain.MigrateRepository = (*migrateRepo)(nil)

type migrateRepo struct {
	*domain.BaseEntRepo
}

func New(dbClient *ent.Client) *migrateRepo {
	bRepo := domain.NewBaseEntRepo(dbClient)
	return &migrateRepo{
		BaseEntRepo: bRepo,
	}
}

func (m *migrateRepo) GetMeilisearch(ctx context.Context) (*ent.AdditionalMigrate, error) {
	result, err := m.GetEntClient(ctx).AdditionalMigrate.
		Query().
		Where(
			additionalmigrate.Name(migrateNameMeilisearch),
		).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return &ent.AdditionalMigrate{
				Name:    migrateNameMeilisearch,
				Version: 0,
			}, nil
		}
	}

	return result, nil
}

func (m *migrateRepo) UpdateMeilisearch(ctx context.Context, version int) error {
	err := m.GetEntClient(ctx).AdditionalMigrate.
		Create().
		SetName(migrateNameMeilisearch).
		SetVersion(version).
		OnConflict(
			sql.ConflictColumns(additionalmigrate.FieldName),
			sql.ResolveWithNewValues(),
		).
		Update(func(set *ent.AdditionalMigrateUpsert) {
			set.SetVersion(version)
		}).
		Exec(ctx)
	if err != nil {
		return hserr.NewInternalError(err, "update meilisearch")
	}

	return nil
}
