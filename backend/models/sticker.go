package models

import (
	"time"
)

type Sticker struct {
	ID               int       `gorm:"column:id;auto_increment;primary_key;not null"`
	StickerName      string    `gorm:"column:stickername;not null"`
	ImgURL           string    `gorm:"column:imgurl;not null"`
	LocalSave        string    `gorm:"column:localsave;not null"`
	IsGIF            bool      `gorm:"column:isgif;not null"` // Deprecated
	LatestUpdateTime time.Time `gorm:"column:latestupdatetime;not null;default:now()"`
}

// TableName 指定 Image 表格的名稱
func (Sticker) TableName() string {
	return "sticker"
}

func StickerCreate(s *Sticker) error {
	res := db.Create(s)
	if res.Error != nil {
		return res.Error
	}

	return nil
}

func StickerGetByID(id int) (*Sticker, bool, error) {
	s := &Sticker{
		ID: id,
	}
	res := db.Find(s)
	if res.Error != nil {
		return nil, false, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, false, nil
	}

	return s, true, nil
}

// StickerRandGetBySN 根據 sticker name 隨機挑選一個 stickre
func StickerRandGetBySN(name string) (*Sticker, bool, error) {
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

func StickerNameList(offset, limit int) ([]string, int64, error) {
	var stickerNameList []string

	var total int64
	query := db.Model(&Sticker{}).Select("stickername").Group("stickername")
	res := query.Count(&total)
	if res.Error != nil {
		return nil, 0, res.Error
	}
	if total == 0 {
		return nil, 0, nil
	}

	res = query.
		Offset(offset).Limit(limit).Order("stickername desc").
		Find(&stickerNameList)
	if res.Error != nil {
		return nil, 0, res.Error
	}

	return stickerNameList, total, nil
}

func StickerListAllByName(name string) ([]*Sticker, error) {
	var stickerList []*Sticker

	res := db.Where("stickername = ?", name).Find(&stickerList)
	if res.Error != nil {
		return nil, res.Error
	}

	return stickerList, nil
}

func StickerDelete(idList []int) error {
	res := db.Delete(&Sticker{}, idList)
	if res.Error != nil {
		return res.Error
	}

	return nil
}
