package delivery

import (
	"backend/app/domain"
	"backend/app/pkg/ginext"

	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"
)

type stickerController struct {
	stickerUsecase domain.StickerUsecase
}

func NewStickerController(e *gin.Engine, stickerUsecase domain.StickerUsecase) {
	ctrl := stickerController{
		stickerUsecase: stickerUsecase,
	}

	e.POST("/sticker-images", ginext.BindHandler(ctrl.AddStickerImage))
	e.GET("/stickers", ginext.BindHandler(ctrl.ListSticker))
	e.DELETE("/stickers/:id", ginext.BindUriHandler(ctrl.DeleteSticker))
}

func (c *stickerController) AddStickerImage(ctx *gin.Context, req *addImageReq) (*ginext.EmptyResp, error) {
	err := c.stickerUsecase.AddImageByURL(ctx, req.StickerName, req.ImageURL)
	if err != nil {
		return nil, xerrors.Errorf("add image by url: %w", err)
	}

	return nil, nil
}

func (c *stickerController) ListSticker(ctx *gin.Context, req listStickerReq) (*listStickerResp, error) {
	offset := (req.Page - 1) * req.Limit

	listOpts := []domain.StickerListOptionFunc{
		domain.StickerListWithImages(4),
	}
	if req.Search != "" {
		listOpts = append(listOpts, domain.StickerListWithSearchName(req.Search))
	}

	stickerReulst, err := c.stickerUsecase.ListStickers(ctx, offset, req.Limit, listOpts...)
	if err != nil {
		return nil, xerrors.Errorf("list stickers: %w", err)
	}

	resp := newlistStickerRespFromListResult(stickerReulst)
	return resp, nil
}

func (c *stickerController) DeleteSticker(ctx *gin.Context, req *deleteStickerReq) (*ginext.EmptyResp, error) {
	err := c.stickerUsecase.Delete(ctx, req.ID)
	if err != nil {
		return nil, xerrors.Errorf("delete sticker: %w", err)
	}

	return nil, nil
}
