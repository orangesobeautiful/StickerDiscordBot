package repository

import (
	"context"

	"backend/app/domain"
	"backend/app/ent"
	"backend/app/ent/discordguild"
	"backend/app/ent/predicate"
	"backend/app/ent/sticker"
	"backend/app/pkg/hserr"

	"github.com/meilisearch/meilisearch-go"
	"golang.org/x/xerrors"
)

var _ domain.StickerRepository = (*stickerRepository)(nil)

type stickerRepository struct {
	*domain.BaseEntRepo

	meilisearch meilisearch.IndexManager
}

func New(
	client *ent.Client,
	meilisearchServiceManager meilisearch.ServiceManager,
) domain.StickerRepository {
	bRepo := domain.NewBaseEntRepo(client)

	meilisearchServiceManager.Index("sticker")

	return &stickerRepository{
		BaseEntRepo: bRepo,
	}
}

func (r *stickerRepository) CreateIfNotExist(ctx context.Context, guildID, name string) (stickerID int, err error) {
	s, err := r.GetEntClient(ctx).Sticker.
		Query().
		Where(
			sticker.And(
				sticker.HasGuildWith(discordguild.ID(guildID)),
				sticker.Name(name),
			),
		).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			stickerID, err = r.create(ctx, guildID, name)
			if err != nil {
				return 0, xerrors.Errorf("create sticker: %w", err)
			}
			return stickerID, nil
		}

		return 0, hserr.NewInternalError(err, "query sticker")
	}

	return s.ID, nil
}

func (r *stickerRepository) create(ctx context.Context, guildID, name string) (int, error) {
	s, err := r.GetEntClient(ctx).Sticker.
		Create().
		SetGuildID(guildID).
		SetName(name).
		Save(ctx)
	if err != nil {
		return 0, hserr.NewInternalError(err, "create sticker")
	}
	return s.ID, nil
}

func (r *stickerRepository) GetStickerWithGuildByID(ctx context.Context, stickerID int) (result *ent.Sticker, err error) {
	s, err := r.GetEntClient(ctx).Sticker.
		Query().
		Where(sticker.ID(stickerID)).
		WithGuild().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NewHsNotFoundError("sticker")
		}

		return nil, hserr.NewInternalError(err, "query sticker")
	}

	return s, nil
}

func (r *stickerRepository) FindByName(ctx context.Context, guildID, name string) (result *ent.Sticker, err error) {
	s, err := r.GetEntClient(ctx).Sticker.
		Query().
		Where(
			sticker.And(
				sticker.HasGuildWith(discordguild.ID(guildID)),
				sticker.Name(name),
			),
		).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}

		return nil, hserr.NewInternalError(err, "query sticker")
	}

	return s, nil
}

func (r *stickerRepository) List(
	ctx context.Context, guildID string, offset, limit int, opts ...domain.StickerListOptionFunc,
) (result domain.ListStickerResult, err error) {
	listOpts := domain.NewStickerListOption(opts...)

	queryFilter := r.GetEntClient(ctx).Sticker.
		Query()

	andConds := []predicate.Sticker{
		sticker.HasGuildWith(discordguild.IDEQ(guildID)),
	}
	searchName := listOpts.GetSearchName()
	if searchName != "" {
		andConds = append(andConds, sticker.NameContainsFold(searchName))
	}
	queryFilter = queryFilter.Where(sticker.And(andConds...))

	withImg := listOpts.GetWithImages()
	if withImg {
		imgLimit := listOpts.GetWithImagesLimit()
		queryFilter = queryFilter.WithImages(func(imgQuery *ent.ImageQuery) {
			if imgLimit != 0 {
				imgQuery.Limit(int(imgLimit) * limit)
			}
		})
	}

	total, err := queryFilter.Clone().Count(ctx)
	if err != nil {
		return result, hserr.NewInternalError(err, "query sticker count")
	}

	queryFilter = queryFilter.
		Offset(offset).
		Limit(limit)
	ss, err := queryFilter.All(ctx)
	if err != nil {
		return result, hserr.NewInternalError(err, "query sticker")
	}

	result = domain.NewListResult(total, ss)
	return result, nil
}

func (r *stickerRepository) AddImage(ctx context.Context, stickerID int, imageIDs ...int) (err error) {
	_, err = r.GetEntClient(ctx).Sticker.
		UpdateOneID(stickerID).
		AddImageIDs(imageIDs...).
		Save(ctx)
	if err != nil {
		return hserr.NewInternalError(err, "add image to sticker")
	}

	return nil
}

func (r *stickerRepository) Delete(ctx context.Context, id ...int) (err error) {
	_, err = r.GetEntClient(ctx).Sticker.
		Delete().
		Where(sticker.IDIn(id...)).
		Exec(ctx)
	if err != nil {
		return hserr.NewInternalError(err, "delete sticker")
	}

	return nil
}
