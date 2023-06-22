package models

import (
	"time"
)

type Sticker struct {
	ID               int       `gorm:"id;auto_increment;primary_key;not null"`
	StickerName      string    `gorm:"stickername;not null"`
	ImgURL           string    `gorm:"imgurl;not null"`
	LocalSave        string    `gorm:"localsave;not null"`
	IsGIF            bool      `gorm:"isgif;not null"`
	LatestUpdateTime time.Time `gorm:"latestupdatetime;not null;default:now()"`
}

// TableName 指定 Image 表格的名稱
func (Sticker) TableName() string {
	return "sticker"
}

// StickerRandFindBySN 根據 sticker name 隨機挑選一個 stickre
func StickerRandFindBySN(name string) (*Sticker, bool, error) {
	s := &Sticker{
		StickerName: name,
	}
	res := db.Where(s).Limit(1).Order("rand()").Find(s)
	if res.Error != nil {

		return nil, false, res.Error
	}

	if res.RowsAffected == 0 {
		return nil, false, nil
	}

	return s, true, nil
}
