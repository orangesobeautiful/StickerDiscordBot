package domainresponse

import (
	"backend/app/ent"
)

type DiscordUser struct {
	ID int `json:"id"`

	DiscordID string `json:"discord_id"`

	GuildID string `json:"guild_id"`

	Name string `json:"name"`

	AvatarURL string `json:"avatar_url"`
}

func (rd *DomainResponse) NewIDiscordUserFromEnt(dcUser *ent.DiscordUser) *DiscordUser {
	return &DiscordUser{
		ID:        dcUser.ID,
		DiscordID: dcUser.DiscordID,
		GuildID:   dcUser.GuildID,
		Name:      dcUser.Name,
		AvatarURL: dcUser.AvatarURL,
	}
}
