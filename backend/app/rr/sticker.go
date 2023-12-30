package rr




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
