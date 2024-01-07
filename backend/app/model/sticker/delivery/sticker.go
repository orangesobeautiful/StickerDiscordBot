package delivery

import (
	"context"

	"backend/app/domain"
	domainresponse "backend/app/domain-response"
	ginauth "backend/app/model/discorduser/gin-auth"
	discordcommand "backend/app/pkg/discord-command"
	"backend/app/pkg/ginext"

	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"
)

type stickerController struct {
	auth           ginauth.AuthInterface
	rd             *domainresponse.DomainResponse
	stickerUsecase domain.StickerUsecase
}

func Initialze(
	apiGroup *gin.RouterGroup, dcCmdRegister discordcommand.Register,
	auth ginauth.AuthInterface,
	rd *domainresponse.DomainResponse,
	stickerUsecase domain.StickerUsecase,
) {
	ctrl := stickerController{
		auth:           auth,
		rd:             rd,
		stickerUsecase: stickerUsecase,
	}

	ctrl.RegisterGinRouter(apiGroup)
	ctrl.RegisterDiscordCommand(dcCmdRegister)
}

func (c *stickerController) AddStickerImage(ctx context.Context, req *addImageReq) (*ginext.EmptyResp, error) {
	err := c.stickerUsecase.AddImageByURL(ctx, req.StickerName, req.ImageURL)
	if err != nil {
		return nil, xerrors.Errorf("add image by url: %w", err)
	}

	return nil, nil
}

func (c *stickerController) ListSticker(ctx context.Context, req listStickerReq) (*listStickerResp, error) {
	offset := (req.Page - 1) * req.Limit

	const maxImagePreviewLimit = 4
	listOpts := []domain.StickerListOptionFunc{
		domain.StickerListWithImages(maxImagePreviewLimit),
	}
	if req.Search != "" {
		listOpts = append(listOpts, domain.StickerListWithSearchName(req.Search))
	}

	stickerReulst, err := c.stickerUsecase.ListStickers(ctx, offset, req.Limit, listOpts...)
	if err != nil {
		return nil, xerrors.Errorf("list stickers: %w", err)
	}

	resp := c.newlistStickerRespFromListResult(stickerReulst)
	return resp, nil
}

func (c *stickerController) DeleteSticker(ctx context.Context, req *deleteStickerReq) (*ginext.EmptyResp, error) {
	err := c.stickerUsecase.Delete(ctx, req.ID)
	if err != nil {
		return nil, xerrors.Errorf("delete sticker: %w", err)
	}

	return nil, nil
}

func (c *stickerController) DeleteStickerByName(ctx context.Context, req *deleteStickerByNameReq) (*ginext.EmptyResp, error) {
	err := c.stickerUsecase.DeleteByName(ctx, req.Name)
	if err != nil {
		return nil, xerrors.Errorf("delete sticker by name: %w", err)
	}

	return nil, nil
}
