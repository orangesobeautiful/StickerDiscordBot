package delivery

import (
	domainresponse "backend/app/domain-response"
	"backend/app/ent"
	discordcommand "backend/app/pkg/discord-command"
)

type createLoginCodeResp struct {
	Code string `json:"code"`
}

type checkLoginCodeReq struct {
	Code string `json:"code" binding:"required"`
}

type checkLoginCodeResp struct {
	IsVerified bool `json:"is_verified"`

	Token string `json:"token"`
}

type verifyLoginCodeReq struct {
	discordcommand.BaseAuthInteractionCreate `dccmd:"ignore"`

	Code string `dccmd:"name=code"`
}

type getSelfInfoResp = domainresponse.DiscordUser

func (c *discorduserController) newGetSelfInfoRespFromEnt(dcUser *ent.DiscordUser) *getSelfInfoResp {
	return c.rd.NewIDiscordUserFromEnt(dcUser)
}
