package controllers

import (
	"time"

	"backend/models"
	"backend/pkg/ginext"
	"backend/pkg/hserr"
	"backend/pkg/log"
	"backend/rr"
	"backend/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (c *Controller) WebGenLoginCode(ctx *gin.Context) (*rr.WebGenLoginCodeResp, error) {
	randCode := utils.RandString([]byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"), 6)
	expDur := 5 * time.Minute
	expTime := time.Now().Add(expDur)
	err := models.WebLoginVerificationCreate(randCode, expTime)
	if err != nil {
		return nil, hserr.ErrInternalServerError
	}

	return &rr.WebGenLoginCodeResp{
		Code: randCode,
	}, nil
}

func (c *Controller) WebCheckLogin(ctx *gin.Context, req *rr.WebCheckLoginReq) (*rr.WebCheckLoginResp, error) {
	info, exist, err := models.WebLoginVerificationGetByCode(req.Code)
	if err != nil {
		return nil, hserr.ErrInternalServerError
	}

	if !exist || info.UserID == "" {
		return &rr.WebCheckLoginResp{Result: 0}, nil
	}
	if time.Now().After(info.Expirationtime) {
		return &rr.WebCheckLoginResp{Result: 0}, nil
	}

	sess := sessions.Default(ctx)
	sess.Set("user-auth", info.UserID)
	if err = sess.Save(); err != nil {
		return nil, hserr.ErrInternalServerError
	}

	return &rr.WebCheckLoginResp{
		Result: 1,
		UserID: info.UserID,
	}, nil
}

func (c *Controller) WebLogout(ctx *gin.Context) (*ginext.EmptyResp, error) {
	sess := sessions.Default(ctx)
	sess.Delete("user-auth")
	if err := sess.Save(); err != nil {
		log.Errorf("sess.Save failed, err=%s", err)
		return nil, hserr.ErrInternalServerError
	}

	return nil, nil
}

func (c *Controller) WebHasLlogin(ctx *gin.Context) (*ginext.EmptyResp, error) {
	return nil, nil
}

func (c *Controller) WebSelfInfo(ctx *gin.Context) (*rr.WebSelfInfoResp, error) {
	sess := sessions.Default(ctx)
	userID := sess.Get("user-auth").(string)
	userInfo, exist, err := models.WebUserInfoGetByID(userID)
	if err != nil {
		log.Errorf("WebUserInfoGetByID failed, err=%s", err)
		return nil, hserr.ErrInternalServerError
	}
	if !exist {
		sess.Delete("user-auth")
		_ = sess.Save()
		return nil, hserr.ErrUnauthorized
	}

	return &rr.WebSelfInfoResp{
		Name:      userInfo.Name,
		AvatarURL: userInfo.AvatarURL,
	}, nil
}
