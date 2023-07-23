package models

import (
	"time"

	"backend/pkg/log"
)

type WebLoginVerification struct {
	Code           string    `gorm:"column:code;type:varchar(255);primary_key"`
	Expirationtime time.Time `gorm:"column:expirationtime;not null"`
	UserID         string    `gorm:"column:userid;type:varchar(255)"`
}

// TableName Specify the name of the WebLoginVerification table
func (WebLoginVerification) TableName() string {
	return "webloginverification"
}

func WebLoginVerificationCreate(code string, expireTime time.Time) error {
	v := &WebLoginVerification{
		Code:           code,
		Expirationtime: expireTime,
	}
	res := db.Create(v)
	if res.Error != nil {
		log.Errorf("WebLoginVerificationCreate failed, err=%s", res.Error)
		return res.Error
	}

	return nil
}

func WebLoginVerificationGetByCode(code string) (*WebLoginVerification, bool, error) {
	v := &WebLoginVerification{
		Code: code,
	}
	res := db.Where(v).Find(v)
	if res.Error != nil {
		log.Errorf("WebLoginVerificationGetByCode failed, err=%s", res.Error)
		return nil, false, res.Error
	}

	if res.RowsAffected == 0 {
		return nil, false, nil
	}

	return v, true, nil
}

type WebUserInfo struct {
	ID        string `gorm:"column:userid;type:varchar(255);primary_key"`
	Name      string `gorm:"column:name;type:varchar(255);not null"`
	AvatarURL string `gorm:"column:avatarurl;type:text;not null"`
}

// TableName Specify the name of the WebUserInfo table
func (WebUserInfo) TableName() string {
	return "webuserinfo"
}

func WebUserInfoGetByID(id string) (*WebUserInfo, bool, error) {
	u := &WebUserInfo{
		ID: id,
	}
	res := db.Where(u).Find(u)
	if res.Error != nil {
		log.Errorf("WebUserInfoGetByID failed, err=%s", res.Error)
		return nil, false, res.Error
	}

	if res.RowsAffected == 0 {
		return nil, false, nil
	}

	return u, true, nil
}
