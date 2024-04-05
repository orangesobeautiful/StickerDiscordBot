package domainresponse

import (
	"backend/app/ent"
)

type Sticker struct {
	ID int `json:"id"`

	StickerName string `json:"sticker_name"`

	Images []*Image `json:"images"`
}

func (rd *DomainResponse) NewStickerFromEnt(s *ent.Sticker) *Sticker {
	imgs := make([]*Image, len(s.Edges.Images))
	for i, entImg := range s.Edges.Images {
		imgs[i] = rd.NewImageFromEnt(entImg)
	}

	return &Sticker{
		ID:          s.ID,
		StickerName: s.Name,
		Images:      imgs,
	}
}

type Image struct {
	ID int `json:"id"`

	URL string `json:"url"`
}

func (rd *DomainResponse) NewImageFromEnt(img *ent.Image) *Image {
	imgURL := rd.objDataConvert.GetObjectDirectURL(img.SavePath)

	return &Image{
		ID:  img.ID,
		URL: imgURL,
	}
}
