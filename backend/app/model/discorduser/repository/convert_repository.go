package repository

import (
	"backend/app/domain"
)

func convertDCWebLoginVerifyInfo(code string, info dcWebLoginVerifyInfo) *domain.DiscordWebLoginVerification {
	return domain.NewDiscordWebLoginLoginVerification(
		code,
		info.DCUserID,
		info.DCChannelID,
	)
}
