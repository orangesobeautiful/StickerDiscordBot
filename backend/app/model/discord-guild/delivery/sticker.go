package delivery

import (
	"context"

	"backend/app/domain"
	"backend/app/pkg/ginext"

	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"
)

func (c *discordGuildController) ginAddStickerImage(ctx *gin.Context, req *ginAddImageReq) (*ginext.EmptyResp, error) {
	user := c.auth.MustGetUserFromContext(ctx)

	err := c.addStickerImage(ctx, user.GuildID, req.StickerName, req.ImageURL)
	if err != nil {
		return nil, xerrors.Errorf("add sticker image: %w", err)
	}

	return nil, nil
}

func (c *discordGuildController) addStickerImage(ctx context.Context, guildID, stickerName, imageURL string) error {
	err := c.stickerUsecase.AddImageByURL(ctx, guildID, stickerName, imageURL)
	if err != nil {
		return xerrors.Errorf("add image by url: %w", err)
	}

	return nil
}

func (c *discordGuildController) ginListSticker(ctx *gin.Context, req ginListStickerReq) (*listStickerResp, error) {
	user := c.auth.MustGetUserFromContext(ctx)

	resp, err := c.listSticker(ctx, user.GuildID, req.Page, req.Limit, req.Search)
	if err != nil {
		return nil, xerrors.Errorf("list sticker: %w", err)
	}

	return resp, nil
}

func (c *discordGuildController) listSticker(
	ctx context.Context, guildID string, page, limit int, searchName string,
) (*listStickerResp, error) {
	offset := (page - 1) * limit

	const maxImagePreviewLimit = 4
	listOpts := []domain.StickerListOptionFunc{
		domain.StickerListWithImages(maxImagePreviewLimit),
	}
	if searchName != "" {
		listOpts = append(listOpts, domain.StickerListWithSearchName(searchName))
	}

	stickerReulst, err := c.stickerUsecase.ListStickers(ctx, guildID, offset, limit, listOpts...)
	if err != nil {
		return nil, xerrors.Errorf("list stickers: %w", err)
	}

	resp := c.newlistStickerRespFromListResult(stickerReulst)
	return resp, nil
}

func (c *discordGuildController) ginDeleteSticker(ctx *gin.Context, req ginDeleteStickerReq) (*ginext.EmptyResp, error) {
	err := c.deleteSticker(ctx, req.StickerID)
	if err != nil {
		return nil, xerrors.Errorf("delete sticker: %w", err)
	}

	return nil, nil
}

func (c *discordGuildController) deleteSticker(ctx context.Context, stickerID int) error {
	err := c.stickerUsecase.Delete(ctx, stickerID)
	if err != nil {
		return xerrors.Errorf("delete sticker: %w", err)
	}

	return nil
}

func (c *discordGuildController) deleteStickerByName(
	ctx context.Context, guildID string, stickerName string,
) error {
	err := c.stickerUsecase.DeleteByName(ctx, guildID, stickerName)
	if err != nil {
		return xerrors.Errorf("delete sticker by name: %w", err)
	}

	return nil
}
