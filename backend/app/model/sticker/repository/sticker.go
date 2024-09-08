package repository

import (
	"context"
	"strconv"

	"backend/app/domain"
	"backend/app/ent"
	"backend/app/ent/discordguild"
	"backend/app/ent/image"
	"backend/app/ent/sticker"
	"backend/app/pkg/hserr"

	"github.com/meilisearch/meilisearch-go"
	"golang.org/x/xerrors"
)

var _ domain.StickerRepository = (*stickerRepository)(nil)

type stickerRepository struct {
	*domain.BaseEntRepo

	meilisearchSticker meilisearch.IndexManager
}

func New(
	client *ent.Client,
	meilisearchServiceManager meilisearch.ServiceManager,
	meilisearchIndexName domain.MeilisearchIndexName,
) domain.StickerRepository {
	bRepo := domain.NewBaseEntRepo(client)

	stickerIndexName := meilisearchIndexName.Sticker()
	stickerIndex := meilisearchServiceManager.Index(stickerIndexName)

	return &stickerRepository{
		BaseEntRepo:        bRepo,
		meilisearchSticker: stickerIndex,
	}
}

func (r *stickerRepository) CreateIfNotExist(ctx context.Context, guildID, name string) (stickerID int, err error) {
	handler := newCreateHandler(
		r,
		r.meilisearchSticker,
		guildID,
		name,
	)

	err = r.WithTx(ctx, handler.Do)
	if err != nil {
		return 0, xerrors.Errorf("create sticker: %w", err)
	}

	return handler.stickerResult.ID, nil
}

type createHandler struct {
	*domain.BaseEntRepo

	meilisearchSticker meilisearch.IndexManager

	guildID string

	name string

	existed bool

	stickerResult *ent.Sticker
}

func newCreateHandler(
	repo *stickerRepository,
	meilisearchSticker meilisearch.IndexManager,
	guildID,
	name string,
) *createHandler {
	return &createHandler{
		BaseEntRepo:        repo.BaseEntRepo,
		meilisearchSticker: meilisearchSticker,
		guildID:            guildID,
		name:               name,
	}
}

func (h *createHandler) Do(ctx context.Context) error {
	err := h.handleBaseDB(ctx)
	if err != nil {
		return xerrors.Errorf("handle base db: %w", err)
	}

	if h.existed {
		return nil
	}

	err = h.handleMeilisearch(ctx)
	if err != nil {
		return xerrors.Errorf("handle meilisearch: %w", err)
	}

	return nil
}

func (h *createHandler) handleBaseDB(ctx context.Context) error {
	s, err := h.GetEntClient(ctx).Sticker.
		Query().
		Where(
			sticker.And(
				sticker.HasGuildWith(discordguild.ID(h.guildID)),
				sticker.Name(h.name),
			),
		).
		Only(ctx)
	if err == nil {
		h.existed = true
	} else {
		if ent.IsNotFound(err) {
			return h.createToBaseDB(ctx)
		}

		return hserr.NewInternalError(err, "query sticker")
	}

	h.stickerResult = s

	return nil
}

func (h *createHandler) createToBaseDB(ctx context.Context) error {
	s, err := h.GetEntClient(ctx).Sticker.
		Create().
		SetGuildID(h.guildID).
		SetName(h.name).
		Save(ctx)
	if err != nil {
		return hserr.NewInternalError(err, "create sticker")
	}

	h.stickerResult = s

	return nil
}

func (h *createHandler) handleMeilisearch(ctx context.Context) error {
	meilisearchSticker := newStickerMeilisearchEntity(
		h.stickerResult.ID,
		h.stickerResult.Name,
		h.guildID,
		h.stickerResult.CreatedAt,
	)

	taskInfo, err := h.meilisearchSticker.AddDocumentsWithContext(ctx, meilisearchSticker)
	if err != nil {
		return hserr.NewInternalError(err, "add sticker to meilisearch")
	}

	task, err := h.meilisearchSticker.WaitForTaskWithContext(ctx, taskInfo.TaskUID, 0)
	if err != nil {
		return hserr.NewInternalError(err, "wait for task")
	}

	if task.Status != meilisearch.TaskStatusSucceeded {
		return hserr.NewInternalError(err, "task failed")
	}

	return nil
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

	query := listOpts.GetSearchName()

	searchResult, err := r.searchWithMeilisearch(ctx, guildID, query, offset, limit)
	if err != nil {
		return result, xerrors.Errorf("search with meilisearch: %w", err)
	}

	matchedIDs := make([]int, len(searchResult.Hits))
	for i, hit := range searchResult.Hits {
		matchedIDs[i] = hit.ID
	}

	queryFilter := r.GetEntClient(ctx).Sticker.
		Query().
		Where(
			sticker.IDIn(matchedIDs...),
		)

	withImg := listOpts.GetWithImages()
	if withImg {
		imgLimit := listOpts.GetWithImagesLimit()
		queryFilter = queryFilter.WithImages(func(imgQuery *ent.ImageQuery) {
			if imgLimit != 0 {
				imgQuery.Limit(int(imgLimit) * limit)
			}
		})
	}

	ss, err := queryFilter.All(ctx)
	if err != nil {
		return result, hserr.NewInternalError(err, "query sticker")
	}

	result = domain.NewListResult(searchResult.EstimatedTotalHits, ss)
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
	handler := newDeleteHandler(
		r,
		r.meilisearchSticker,
		id,
	)

	err = r.WithTx(ctx, handler.Do)
	if err != nil {
		return xerrors.Errorf("delete sticker: %w", err)
	}

	return nil
}

type deleteHandler struct {
	*domain.BaseEntRepo

	meilisearchSticker meilisearch.IndexManager

	stickerIDs []int
}

func newDeleteHandler(
	repo *stickerRepository,
	meilisearchSticker meilisearch.IndexManager,
	stickerIDs []int,
) *deleteHandler {
	return &deleteHandler{
		BaseEntRepo:        repo.BaseEntRepo,
		meilisearchSticker: meilisearchSticker,
		stickerIDs:         stickerIDs,
	}
}

func (h *deleteHandler) Do(ctx context.Context) error {
	err := h.handleBaseDB(ctx)
	if err != nil {
		return xerrors.Errorf("handle base db: %w", err)
	}

	err = h.handleMeilisearch(ctx)
	if err != nil {
		return xerrors.Errorf("handle meilisearch: %w", err)
	}

	return nil
}

func (h *deleteHandler) handleBaseDB(ctx context.Context) error {
	_, err := h.GetEntClient(ctx).Sticker.
		Delete().
		Where(sticker.IDIn(h.stickerIDs...)).
		Exec(ctx)
	if err != nil {
		return hserr.NewInternalError(err, "delete sticker")
	}

	return nil
}

func (h *deleteHandler) handleMeilisearch(ctx context.Context) error {
	identifiers := make([]string, 0, len(h.stickerIDs))

	for _, id := range h.stickerIDs {
		identifiers = append(identifiers, strconv.Itoa(id))
	}

	taskInfo, err := h.meilisearchSticker.DeleteDocumentsWithContext(ctx, identifiers)
	if err != nil {
		return hserr.NewInternalError(err, "delete sticker from meilisearch")
	}

	task, err := h.meilisearchSticker.WaitForTaskWithContext(ctx, taskInfo.TaskUID, 0)
	if err != nil {
		return hserr.NewInternalError(err, "wait for task")
	}

	if task.Status != meilisearch.TaskStatusSucceeded {
		return hserr.NewInternalError(err, "task failed")
	}

	return nil
}

func (r *stickerRepository) GetStickerImageWithGuildByID(ctx context.Context, imageID int) (*ent.Image, error) {
	i, err := r.GetEntClient(ctx).Image.
		Query().
		Where(
			image.IDEQ(imageID),
		).
		WithSticker(
			func(sq *ent.StickerQuery) {
				sq.WithGuild()
			},
		).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, domain.NewHsNotFoundError("sticker image")
		}

		return nil, hserr.NewInternalError(err, "query sticker image")
	}

	return i, nil
}
