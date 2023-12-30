package delivery

import (
	"backend/app/domain"
	"backend/app/ent"
)

type addImageReq struct {
	StickerName string `json:"sticker_name" binding:"required"`
	ImageURL    string `json:"image_url" binding:"required,http_url"`
}

type listStickerReq struct {
	Page int `form:"page" binding:"required,gte=1"`

	Limit int `form:"limit" binding:"required,gte=1,lte=30"`

	Search string `form:"search"`
}

type listStickerResp struct {
	TotalCount int `json:"total_count"`

	Stickers []*sticker `json:"stickers"`
}

func newlistStickerRespFromListResult(listResult domain.ListStickerResult) *listStickerResp {
	entSs := listResult.GetItems()
	ss := make([]*sticker, len(entSs))
	for i, entS := range entSs {
		ss[i] = newStickerEntSticker(entS)
	}

	return &listStickerResp{
		TotalCount: listResult.GetTotal(),
		Stickers:   ss,
	}
}

type sticker struct {
	ID int `json:"id"`

	StickerName string `json:"sticker_name"`

	Images []*image `json:"images"`
}

func newStickerEntSticker(s *ent.Sticker) *sticker {
	imgs := make([]*image, len(s.Edges.Images))
	for i, entImg := range s.Edges.Images {
		imgs[i] = newImageFromEntImage(entImg)
	}

	return &sticker{
		ID:          s.ID,
		StickerName: s.Name,
		Images:      imgs,
	}
}

type image struct {
	ID int `json:"id"`

	URL string `json:"url"`
}

func newImageFromEntImage(img *ent.Image) *image {
	return &image{
		ID:  img.ID,
		URL: img.SavePath,
	}
}

type deleteStickerReq struct {
	ID int `uri:"id" binding:"required,gte=0"`
}
