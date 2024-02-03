package server

import (
	domainresponse "backend/app/domain-response"
	chatdelivery "backend/app/model/chat/delivery"
	chatrepository "backend/app/model/chat/repository"
	chatusecase "backend/app/model/chat/usecase"
	debugdelivery "backend/app/model/debug/delivery"
	discordguildrepository "backend/app/model/discord-guild/repository"
	discordguildusecase "backend/app/model/discord-guild/usecase"
	discordmessage "backend/app/model/discord-message"
	discorduserdelivery "backend/app/model/discorduser/delivery"
	ginauth "backend/app/model/discorduser/gin-auth"
	discorduserrepo "backend/app/model/discorduser/repository"
	discorduserusecase "backend/app/model/discorduser/usecase"
	imagerepo "backend/app/model/image/repository"
	imageusecase "backend/app/model/image/usecase"
	stickerdelivery "backend/app/model/sticker/delivery"
	stickerrepo "backend/app/model/sticker/repository"
	stickerusecase "backend/app/model/sticker/usecase"
	discordcommand "backend/app/pkg/discord-command"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"golang.org/x/xerrors"
)

func (s *Server) setModel(
	sessStore sessions.Store,
	e *gin.Engine, dcCmdRegister discordcommand.Register,
	rd *domainresponse.DomainResponse,
) (dcMsgHandler discordmessage.HandlerInterface, err error) {
	apiGroup := e.Group("/api/v1")

	dcUserRepo := discorduserrepo.NewDiscordUser(s.dbClient)
	dcUserWebLoginRepo := discorduserrepo.NewRedisDCWebLoginVerification(s.dbClient, s.redisClient)
	dcUserWebLoginUsecase := discorduserusecase.NewDCWebUsecase(dcUserWebLoginRepo, dcUserRepo)
	auth := ginauth.New(sessStore, dcUserWebLoginUsecase, ginauth.WithErrRespHandler(ginHSERROutput))
	discorduserdelivery.Initialze(apiGroup, dcCmdRegister, auth, rd, dcUserWebLoginUsecase)

	debugdelivery.Initialze(apiGroup, auth)

	discordGuildRepo := discordguildrepository.New(s.dbClient)
	discordGuildUsecase := discordguildusecase.New(discordGuildRepo)

	chatRepo := chatrepository.New(s.dbClient)
	chatUsecase := chatusecase.New(chatRepo, discordGuildUsecase, s.openaiCli)
	chatdelivery.Initialze(apiGroup, dcCmdRegister, auth, rd, chatUsecase)

	imageRepo, err := imagerepo.New(s.dbClient, s.bucketHandler)
	if err != nil {
		return nil, xerrors.Errorf("new image repo: %w", err)
	}
	imageUsecase := imageusecase.New(imageRepo)
	_ = imageUsecase

	stickerRepo := stickerrepo.New(s.dbClient)
	stickerUsecase := stickerusecase.New(stickerRepo, imageRepo)
	stickerdelivery.Initialze(apiGroup, dcCmdRegister, auth, rd, stickerUsecase)

	dcMsgHandler = discordmessage.New(stickerUsecase, rd)
	return dcMsgHandler, nil
}
