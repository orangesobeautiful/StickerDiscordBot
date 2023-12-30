package server

import (
	"context"

	"backend/app/config"
	discorduserrepo "backend/app/model/discorduser/repository"
	discorduserusecase "backend/app/model/discorduser/usecase"
	imagerepo "backend/app/model/image/repository"
	imageusecase "backend/app/model/image/usecase"
	stickerdelivery "backend/app/model/sticker/delivery"
	stickerrepo "backend/app/model/sticker/repository"
	stickerusecase "backend/app/model/sticker/usecase"

	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"
)

func (s *Server) setModel(ctx context.Context, e *gin.Engine, cfg *config.CfgInfo) (err error) {
	dcUserRepo := discorduserrepo.NewDiscordUser(s.dbClient)
	dcUserWebLoginRepo := discorduserrepo.NewRedisDCWebLoginVerification(s.redisClient)
	dcUserUsecase := discorduserusecase.NewDiscordUser(dcUserRepo)
	dcUserWebLoginUsecase := discorduserusecase.NewDCWebUsecase(dcUserWebLoginRepo)
	_ = dcUserUsecase
	_ = dcUserWebLoginUsecase

	imageRepo, err := imagerepo.New(ctx, cfg, s.dbClient)
	if err != nil {
		return xerrors.Errorf("new image repo: %w", err)
	}
	imageUsecase := imageusecase.New(imageRepo)
	_ = imageUsecase

	stickerRepo := stickerrepo.New(s.dbClient)
	stickerUsecase := stickerusecase.New(stickerRepo, imageRepo)
	stickerdelivery.NewStickerController(e, stickerUsecase)

	return nil
}
