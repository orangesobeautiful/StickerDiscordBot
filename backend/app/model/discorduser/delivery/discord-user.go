package delivery

import (
	"context"

	"backend/app/domain"
	domainresponse "backend/app/domain-response"
	ginauth "backend/app/model/discorduser/gin-auth"
	discordcommand "backend/app/pkg/discord-command"
	"backend/app/pkg/ginext"

	"github.com/gin-gonic/gin"
	"golang.org/x/xerrors"
)

type discorduserController struct {
	auth ginauth.AuthInterface
	rd   *domainresponse.DomainResponse

	dcWebLoginUsecase domain.DiscordWebLoginVerificationUsecase
}

func Initialze(
	e *gin.Engine, dcCmdRegister discordcommand.Register,
	auth ginauth.AuthInterface,
	rd *domainresponse.DomainResponse,
	dcWebLoginUsecase domain.DiscordWebLoginVerificationUsecase,
) {
	ctrl := discorduserController{
		auth:              auth,
		rd:                rd,
		dcWebLoginUsecase: dcWebLoginUsecase,
	}

	ctrl.RegisterGinRouter(e)
	ctrl.RegisterDiscordCommand(dcCmdRegister)
}

func (c *discorduserController) CreateLoginCode(ctx context.Context) (*createLoginCodeResp, error) {
	code, err := c.dcWebLoginUsecase.CreateRandomLoginCode(ctx)
	if err != nil {
		return nil, xerrors.Errorf("create random login code: %w", err)
	}

	return &createLoginCodeResp{
		Code: code,
	}, nil
}

func (c *discorduserController) VerifyLoginCode(req *verifyLoginCodeReq) (*ginext.EmptyResp, error) {
	err := c.dcWebLoginUsecase.VerifyLoginCode(
		context.Background(),
		req.Code, req.UserID, req.GuildID, req.Name, req.AvatarURL,
	)
	if err != nil {
		return nil, xerrors.Errorf("verify login code: %w", err)
	}

	return nil, nil
}

func (c *discorduserController) CheckLoginCode(ctx *gin.Context, req checkLoginCodeReq) (*ginext.EmptyResp, error) {
	newSessionID, err := c.dcWebLoginUsecase.CheckLoginCode(ctx, req.Code)
	if err != nil {
		return nil, xerrors.Errorf("check login code: %w", err)
	}

	err = c.auth.SetSession(ctx, newSessionID)
	if err != nil {
		return nil, xerrors.Errorf("save session: %w", err)
	}

	return nil, nil
}

func (c *discorduserController) GetSelfInfo(ctx *gin.Context) (*getSelfInfoResp, error) {
	dcUser := c.auth.MustGetUserFromContext(ctx)

	return c.newGetSelfInfoRespFromEnt(dcUser), nil
}
