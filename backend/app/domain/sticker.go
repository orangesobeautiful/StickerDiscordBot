package domain

import (
	"context"

	"backend/app/ent"
)

type ListStickerResult = ListResult[*ent.Sticker]

type StickerRepository interface {
	BaseEntRepoInterface
	CreateIfNotExist(ctx context.Context, name string) (stickerID int, err error)
	FindByName(ctx context.Context, name string) (sticker *ent.Sticker, err error)
	List(ctx context.Context, offset, limit int, opts ...StickerListOptionFunc) (stickers ListStickerResult, err error)
	AddImage(ctx context.Context, stickerID int, imageIDs ...int) (err error)
	Delete(ctx context.Context, id ...int) (err error)
}

type StickerUsecase interface {
	AddImageByURL(ctx context.Context, name string, imageURL string) (err error)
	RandSelectImage(ctx context.Context, stickerName string) (result *ent.Image, err error)

	ListStickers(ctx context.Context, offset, limit int, opts ...StickerListOptionFunc) (stickers ListStickerResult, err error)
	GetStickerAllImages(ctx context.Context, stickerName string) (result []*ent.Image, err error)

	Delete(ctx context.Context, id ...int) (err error)
}
