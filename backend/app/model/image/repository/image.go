package repository

import (
	"context"

	"backend/app/domain"
	"backend/app/ent"
	"backend/app/ent/image"
	"backend/app/ent/sticker"
	"backend/app/pkg/hserr"
	objectstorage "backend/app/pkg/object-storage"

	"golang.org/x/xerrors"
)

var _ domain.ImageRepository = (*imageRepository)(nil)

type imageRepository struct {
	*domain.BaseEntRepo
	objectOperator objectstorage.BucketObjectOperator
}

func New(client *ent.Client, objectOperator objectstorage.BucketObjectOperator) (repo domain.ImageRepository, err error) {
	bRepo := domain.NewBaseEntRepo(client)

	repo = &imageRepository{
		BaseEntRepo:    bRepo,
		objectOperator: objectOperator,
	}

	return repo, nil
}

func (r *imageRepository) CreateWithURL(ctx context.Context, stickerID int, imageURL string) (id int, err error) {
	err = r.WithTx(ctx, func(ctx context.Context) error {
		var txErr error
		newImgID, txErr := r.create(ctx, stickerID)
		if txErr != nil {
			return xerrors.Errorf("create image: %w", txErr)
		}

		result, txErr := r.downloadAndUploadToObjectStorage(ctx, newImgID, imageURL)
		if txErr != nil {
			return xerrors.Errorf("download and upload to object storage: %w", txErr)
		}

		txErr = r.update(ctx, newImgID, result.saveType, result.uploadKey, result.sha256Checksum)
		if txErr != nil {
			return xerrors.Errorf("update image: %w", txErr)
		}

		id = newImgID
		return nil
	})
	if err != nil {
		return 0, xerrors.Errorf(": %w", err)
	}

	return id, nil
}

func (r *imageRepository) create(
	ctx context.Context,
	stickerID int,
) (int, error) {
	i, err := r.GetEntClient(ctx).Image.
		Create().
		SetSaveType(int(domain.ImgSaveTypeNone)).
		SetSavePath("").
		SetSha256Checksum(nil).
		SetStickerID(stickerID).
		Save(ctx)
	if err != nil {
		return 0, hserr.NewInternalError(err, "create image")
	}

	return i.ID, nil
}

func (r *imageRepository) GetBatch(ctx context.Context, ids ...int) (result []*ent.Image, err error) {
	result, err = r.GetEntClient(ctx).Image.
		Query().
		Where(image.IDIn(ids...)).
		All(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, hserr.NewInternalError(err, "query image")
	}

	return result, nil
}

func (r *imageRepository) update(
	ctx context.Context,
	id int,
	saveType domain.ImgSaveType,
	uploadKey string,
	checksum []byte,
) (err error) {
	_, err = r.GetEntClient(ctx).Image.
		UpdateOneID(id).
		SetSaveType(int(saveType)).
		SetSavePath(uploadKey).
		SetSha256Checksum(checksum).
		Save(ctx)
	if err != nil {
		return hserr.NewInternalError(err, "update image")
	}

	return nil
}

func (r *imageRepository) ListAllByStickerID(ctx context.Context, stickerID int) (result []*ent.Image, err error) {
	imgs, err := r.GetEntClient(ctx).Sticker.
		Query().
		Where(sticker.ID(stickerID)).
		QueryImages().All(ctx)
	if err != nil {
		return result, hserr.NewInternalError(err, "query images")
	}

	return imgs, nil
}

func (r *imageRepository) DeleteByImageEnt(ctx context.Context, images ...*ent.Image) (err error) {
	var ids []int
	for _, img := range images {
		ids = append(ids, img.ID)
	}
	_, err = r.GetEntClient(ctx).Image.
		Delete().
		Where(image.IDIn(ids...)).
		Exec(ctx)
	if err != nil {
		return hserr.NewInternalError(err, "delete images")
	}

	var uploadKeys []string
	for _, img := range images {
		uploadKeys = append(uploadKeys, img.SavePath)
	}
	err = r.objectOperator.DeleteObjects(ctx, uploadKeys...)
	if err != nil {
		return hserr.NewInternalError(err, "delete objects")
	}

	return nil
}
