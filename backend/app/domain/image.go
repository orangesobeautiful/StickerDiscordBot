package domain

import (
	"context"

	"backend/app/ent"

	"golang.org/x/exp/constraints"
	"golang.org/x/xerrors"
)

type ImgSaveType int

const (
	ImgSaveTypeNone ImgSaveType = iota

	ImgSaveTypeCloudfare
)

func NewImgSaveType[T constraints.Unsigned](i T) (result ImgSaveType, err error) {
	result = ImgSaveType(i)
	switch result {
	case ImgSaveTypeCloudfare:
		return result, nil
	default:
		return ImgSaveTypeNone, xerrors.Errorf("invalid ImgSaveType : %d", i)
	}
}

type ListImageResult = ListResult[*ent.Image]

type ImageRepository interface {
	BaseEntRepoInterface
	CreateWithURL(ctx context.Context, stickerID int, imageURL string) (id int, err error)
	GetBatch(ctx context.Context, id ...int) (result []*ent.Image, err error)
	ListAllByStickerID(ctx context.Context, stickerID int) (result []*ent.Image, err error)
	DeleteByImageEnt(ctx context.Context, images ...*ent.Image) (err error)
}

type ImageUsecase interface {
	Delete(ctx context.Context, id ...int) (err error)
}
