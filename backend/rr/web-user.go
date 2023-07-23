package rr

type WebGenLoginCodeResp struct {
	Code string `json:"code"` // login auth code
}

type WebCheckLoginReq struct {
	Code string `form:"code" binding:"required"` // login auth code
}

type WebCheckLoginResp struct {
	Result int    `json:"result"`  // 0: fail, 1: success
	UserID string `json:"user_id"` // web user id
}

type WebSelfInfoResp struct {
	Name      string `json:"name"`       // web user name
	AvatarURL string `json:"avatar_url"` // web user avatar url
}
