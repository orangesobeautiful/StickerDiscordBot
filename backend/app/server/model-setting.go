package server

import (
	domainresponse "backend/app/domain-response"
	discordmessage "backend/app/model/discord-message"
	discorduserrepo "backend/app/model/discorduser/repository"
	discorduserusecase "backend/app/model/discorduser/usecase"
	imagerepo "backend/app/model/image/repository"
	imageusecase "backend/app/model/image/usecase"
	stickerdelivery "backend/app/model/sticker/delivery"
	stickerrepo "backend/app/model/sticker/repository"
	stickerusecase "backend/app/model/sticker/usecase"
	discordcommand "backend/app/pkg/discord-command"

	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"
)

func (s *Server) setModel(
	e *gin.Engine, dcCmdRegister discordcommand.Register, rd *domainresponse.DomainResponse,
) (dcMsgHandler discordmessage.HandlerInterface, err error) {
	dcUserRepo := discorduserrepo.NewDiscordUser(s.dbClient)
	dcUserWebLoginRepo := discorduserrepo.NewRedisDCWebLoginVerification(s.redisClient)
	dcUserUsecase := discorduserusecase.NewDiscordUser(dcUserRepo)
	dcUserWebLoginUsecase := discorduserusecase.NewDCWebUsecase(dcUserWebLoginRepo)
	_ = dcUserUsecase
	_ = dcUserWebLoginUsecase

	imageRepo, err := imagerepo.New(s.dbClient, s.bucketHandler)
	if err != nil {
		return nil, xerrors.Errorf("new image repo: %w", err)
	}
	imageUsecase := imageusecase.New(imageRepo)
	_ = imageUsecase

	stickerRepo := stickerrepo.New(s.dbClient)
	stickerUsecase := stickerusecase.New(stickerRepo, imageRepo)
	stickerdelivery.Initialze(e, dcCmdRegister, stickerUsecase, rd)

	dcMsgHandler = discordmessage.New(stickerUsecase, rd)
	return dcMsgHandler, nil
}
