package domain

import (
	"context"

	"backend/app/ent"
)

type ListStickerResult = ListResult[*ent.Sticker]

type StickerRepository interface {
	BaseEntRepoInterface
	CreateIfNotExist(ctx context.Context, guildID string, name string) (stickerID int, err error)
	GetStickerWithGuildByID(ctx context.Context, stickerID int) (sticker *ent.Sticker, err error)
	FindByName(ctx context.Context, guildID, name string) (sticker *ent.Sticker, err error)
	List(ctx context.Context, guildID string, offset, limit int, opts ...StickerListOptionFunc) (stickers ListStickerResult, err error)
	AddImage(ctx context.Context, stickerID int, imageIDs ...int) (err error)
	Delete(ctx context.Context, id ...int) (err error)
}

type StickerUsecase interface {
	AddImageByURL(ctx context.Context, guildID, name string, imageURL string) (err error)
	RandSelectImage(ctx context.Context, guildID, stickerName string) (result *ent.Image, err error)

	FindByName(ctx context.Context, guildID, name string) (sticker *ent.Sticker, err error)
	ListStickers(ctx context.Context, guildID string, offset, limit int, opts ...StickerListOptionFunc) (stickers ListStickerResult, err error)
	GetStickerAllImages(ctx context.Context, guildID, stickerName string) (result []*ent.Image, err error)

	Delete(ctx context.Context, id ...int) (err error)
	DeleteByName(ctx context.Context, guildID, name string) (err error)
}
