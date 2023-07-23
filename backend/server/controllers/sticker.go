package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"backend/models"
	"backend/pkg/hserr"
	"backend/pkg/log"
	"backend/rr"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/h2non/filetype"
)

func (c *Controller) ListSticker(ctx *gin.Context, req *rr.ListStickerReq) (*rr.ListStickerResp, error) {
	stickerNameList, total, err := models.StickerNameList(*req.Start, *req.Num)
	if err != nil {
		log.Errorf("models.StickerNameList failed, err=%s", err)
		return nil, hserr.ErrInternalServerError
	}

	resp := &rr.ListStickerResp{
		MaxPage: total / int64(*req.Num),
	}

	for snIdx, stickerName := range stickerNameList {
		stickerList, err := models.StickerListAllByName(stickerName)
		if err != nil {
			log.Errorf("models.StickerListAllByName failed, err=%s", err)
			return nil, hserr.ErrInternalServerError
		}

		resp.DataList = append(resp.DataList, &rr.ListStickerImageData{
			StickerName: stickerName,
			StickerList: []*rr.ListStickerImageDataSticker{},
		})

		for _, sticker := range stickerList {
			if sticker.LocalSave != "" {
				sticker.ImgURL, _ = url.JoinPath(c.ImgURL, sticker.LocalSave)
			}
			resp.DataList[snIdx].StickerList = append(

				resp.DataList[snIdx].StickerList, &rr.ListStickerImageDataSticker{
					ID:  sticker.ID,
					URL: sticker.ImgURL,
					GIF: sticker.IsGIF,
				})
		}
	}

	return resp, nil
}

func (c *Controller) SearchSticker(ctx *gin.Context, req *rr.SearchStickerReq) (*rr.ListStickerResp, error) {
	resp := &rr.ListStickerResp{
		MaxPage: 1,
	}

	stickerName := req.Query
	stickerList, err := models.StickerListAllByName(stickerName)
	if err != nil {
		log.Errorf("models.StickerListAllByName failed, err=%s", err)
		return nil, hserr.ErrInternalServerError
	}
	if len(stickerList) == 0 {
		return resp, nil
	}

	resp.DataList = append(resp.DataList, &rr.ListStickerImageData{
		StickerName: stickerName,
		StickerList: []*rr.ListStickerImageDataSticker{},
	})

	for _, sticker := range stickerList {
		if sticker.LocalSave != "" {
			sticker.ImgURL, _ = url.JoinPath(c.ImgURL, sticker.LocalSave)
		}
		resp.DataList[0].StickerList = append(
			resp.DataList[0].StickerList, &rr.ListStickerImageDataSticker{
				ID:  sticker.ID,
				URL: sticker.ImgURL,
				GIF: sticker.IsGIF,
			})
	}

	return resp, nil
}

type LimitedBuffer struct {
	buf  []byte
	size int
}

func NewLimitedBuffer(size int) *LimitedBuffer {
	return &LimitedBuffer{
		buf:  make([]byte, 0, size),
		size: size,
	}
}

func (b *LimitedBuffer) Write(p []byte) (n int, err error) {
	orgPLen := len(p)
	if len(b.buf) >= b.size {
		return orgPLen, nil
	}
	if len(b.buf)+len(p) > b.size {
		p = p[:b.size-len(b.buf)]
	}

	b.buf = append(b.buf, p...)
	return orgPLen, nil
}

func (b *LimitedBuffer) Bytes() []byte {
	return b.buf
}

func (c *Controller) ChangeSticker(ctx *gin.Context, req *rr.ChangeStickerReq) (*rr.ChangeStickerResp, error) {
	var err error
	for _, addInfo := range req.Add {
		hc := http.Client{}
		imgResp, err := hc.Get(addInfo.URL)
		if imgResp.StatusCode != 200 {
			log.Errorf("hc.Get failed, err=%s", err)
			return nil, hserr.ErrInternalServerError
		}
		defer imgResp.Body.Close()

		tmpImageName := uuid.New().String()
		imgPath := filepath.Join("sticker-image", tmpImageName)
		imgFile, err := os.Create(imgPath)
		if err != nil {
			log.Errorf("os.Create failed, err=%s", err)
			return nil, hserr.ErrInternalServerError
		}
		defer imgFile.Close()
		md5Hash := md5.New()
		limitBuf := NewLimitedBuffer(3 * 1024)
		writer := io.MultiWriter(imgFile, md5Hash, limitBuf)

		_, err = io.Copy(writer, imgResp.Body)
		if err != nil {
			log.Errorf("io.Copy failed, err=%s", err)
			return nil, hserr.ErrInternalServerError
		}

		fileKind, _ := filetype.Image(limitBuf.Bytes())
		if fileKind == filetype.Unknown {
			return nil, hserr.ErrBadRequest
		}

		err = imgFile.Close()
		if err != nil {
			log.Errorf("imgFile.Close failed, err=%s", err)
			return nil, hserr.ErrInternalServerError
		}
		fileMD5 := hex.EncodeToString(md5Hash.Sum(nil))
		saveFileName := fileMD5 + "." + fileKind.Extension
		err = os.Rename(filepath.Join("sticker-image", tmpImageName),
			filepath.Join("sticker-image", saveFileName))
		if err != nil {
			log.Errorf("os.Rename failed, err=%s", err)
			return nil, hserr.ErrInternalServerError
		}

		err = models.StickerCreate(&models.Sticker{
			StickerName: req.StickerName,
			ImgURL:      addInfo.URL,
			LocalSave:   saveFileName,
		})

		if err != nil {
			log.Errorf("models.StickerCreate failed, err=%s", err)
			return nil, hserr.ErrInternalServerError
		}
	}
	for _, delID := range req.Delete {
		sticker, exist, err := models.StickerGetByID(delID)
		if err != nil {
			log.Errorf("models.StickerGetByID failed, err=%s", err)
			return nil, hserr.ErrInternalServerError
		}
		if !exist {
			continue
		}

		err = models.StickerDelete([]int{sticker.ID})
		if err != nil {
			log.Errorf("models.StickerDeleteByID failed, err=%s", err)
			return nil, hserr.ErrInternalServerError
		}

		if sticker.LocalSave != "" {
			_ = os.Remove(filepath.Join("sticker-image", sticker.LocalSave))
		}
	}
	resp := new(rr.ChangeStickerResp)

	stickerName := req.StickerName
	stickerList, err := models.StickerListAllByName(stickerName)
	if err != nil {
		log.Errorf("models.StickerListAllByName failed, err=%s", err)
		return nil, hserr.ErrInternalServerError
	}
	if len(stickerList) == 0 {
		return resp, nil
	}

	for _, sticker := range stickerList {
		if sticker.LocalSave != "" {
			sticker.ImgURL, _ = url.JoinPath(c.ImgURL, sticker.LocalSave)
		}
		resp.Imgs = append(resp.Imgs,
			rr.NewChangeStickerImgData(
				sticker.ID,
				sticker.ImgURL,
				sticker.IsGIF))
	}

	return resp, nil
}
