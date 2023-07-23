package models

import (
	"strconv"

	"backend/models/sn"
	"backend/pkg/log"
)

func initSetting() error {
	defaultSetting := map[sn.SettingName]string{
		sn.BotPrefix:            "$",
		sn.StickerDownloadCount: strconv.FormatInt(0, 10),
	}
	for name, defaultV := range defaultSetting {
		info := new(BotInfo)
		exist, err := info.GetByName(name)
		if err != nil {
			log.Errorf("info.GetByName failed, err=%s", err)
			return err
		}
		if !exist {
			info.Name = name
			err = info.Update(defaultV)
			if err != nil {
				log.Errorf("info.Update failed, err=%s", err)
				return err
			}
		}
	}

	return nil
}

type BotInfo struct {
	Name  sn.SettingName `gorm:"column:name;primary_key"`
	Value string         `gorm:"column:value;not null"`
}

// TableName 指定 Image 表格的名稱
func (BotInfo) TableName() string {
	return "botinfo"
}

func (info *BotInfo) GetByName(name sn.SettingName) (bool, error) {
	res := db.Where(BotInfo{Name: name}).Find(info)
	if res.Error != nil {
		return false, res.Error
	}
	if res.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

func (info *BotInfo) Update(value string) error {
	res := db.Where(info).Updates(BotInfo{Value: value})
	if res.Error != nil {
		return res.Error
	}
	info.Value = value

	return nil
}

// BotPrefixGet 回傳 BotPrefix 的設定值
func BotPrefixGet() (string, error) {
	info := new(BotInfo)
	if _, err := info.GetByName(sn.BotPrefix); err != nil {
		return "", err
	}

	return info.Value, nil
}

// BotPrefixSet 設定 BotPrefix 的設定值
func BotPrefixSet(value string) error {
	info := new(BotInfo)
	info.Name = sn.BotPrefix
	if _, err := info.GetByName(sn.BotPrefix); err != nil {
		return err
	}

	return nil
}
