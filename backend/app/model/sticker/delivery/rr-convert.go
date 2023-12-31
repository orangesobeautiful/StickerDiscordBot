package delivery

import (
	"backend/app/domain"
	domainresponse "backend/app/domain-response"
)

type addImageReq struct {
	StickerName string `json:"sticker_name" binding:"required"`

	ImageURL string `json:"image_url" binding:"required,http_url"`
}

type listStickerReq struct {
	Page int `form:"page" binding:"required,gte=1"`

	Limit int `form:"limit" binding:"required,gte=1,lte=30"`

	Search string `form:"search"`
}

type listStickerResp struct {
	TotalCount int `json:"total_count"`

	Stickers []*domainresponse.Sticker `json:"stickers"`
}

func (c *stickerController) newlistStickerRespFromListResult(listResult domain.ListStickerResult) *listStickerResp {
	entSs := listResult.GetItems()
	ss := make([]*domainresponse.Sticker, len(entSs))
	for i, entS := range entSs {
		ss[i] = c.rd.NewStickerFromEnt(entS)
	}

	return &listStickerResp{
		TotalCount: listResult.GetTotal(),
		Stickers:   ss,
	}
}

type deleteStickerReq struct {
	ID int `uri:"id" binding:"required,gte=0"`
}
