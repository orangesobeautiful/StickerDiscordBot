package controllers

import (
	"backend/models"
	"backend/pkg/hs"
	"fmt"
	"net/http"
)

type Req struct {
	StickerName string `query:"stickerName" json:"-" validate:"required"`
	Num         int    `query:"num" json:"-" validate:"required"`
	NumPtr      *int   `query:"numPtr" json:"-" validate:"required"`
}

type Resp struct {
	RspMsg string
}

func GetStickerRand(ctx *hs.Context, req *Req) (*models.Sticker, *hs.ErrResp) {
	fmt.Println("req", req)
	if req.NumPtr == nil {
		fmt.Println("req.NumPtr is nil")
	} else {
		fmt.Println(*req.NumPtr)
	}

	sticker, exist, err := models.StickerRandFindBySN(req.StickerName)
	if err != nil {
		return nil, hs.ErrInternalServerError
	}

	if !exist {
		return nil, hs.NewErr(http.StatusNotFound, "not exist")
	}

	return sticker, nil
}
