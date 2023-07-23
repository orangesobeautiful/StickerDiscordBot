package rr

type ListStickerReq struct {
	Start *int `form:"start" binding:"required,gte=0"` // start index
	Num   *int `form:"num" binding:"required,gt=0"`    // number of stickers
}

type ListStickerImageDataSticker struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
	GIF bool   `json:"gif"`
}

type ListStickerImageData struct {
	StickerName string                         `json:"sn"`
	StickerList []*ListStickerImageDataSticker `json:"sts"`
}

type ListStickerResp struct {
	DataList []*ListStickerImageData `json:"img_data"`
	MaxPage  int64                   `json:"maxp"`
}

type SearchStickerReq struct {
	Query string `form:"q" binding:"required"` // sticker name
}

type ChangeStickerReq struct {
	StickerName string `json:"sn" binding:"required"` // sticker name
	Add         []*struct {
		URL string `json:"url" binding:"required"` // sticker url
	} `json:"add"`
	Delete []int `json:"delete"` // sticker id list
}

type ChangeStickerImgData []any

func NewChangeStickerImgData(id int, url string, gif bool) *ChangeStickerImgData {
	return &ChangeStickerImgData{id, url, gif}
}

type ChangeStickerResp struct {
	Imgs []*ChangeStickerImgData `json:"imgs"`
	Err  string                  `json:"err"` // error message
}
