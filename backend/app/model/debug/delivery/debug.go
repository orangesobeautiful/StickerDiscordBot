package delivery

import (
	ginauth "backend/app/model/discorduser/gin-auth"

	"github.com/gin-gonic/gin"
)

type debugController struct {
	auth ginauth.AuthInterface
}

func Initialze(
	apiGroup *gin.RouterGroup,
	auth ginauth.AuthInterface,
) {
	ctrl := debugController{
		auth: auth,
	}

	ctrl.RegisterGinRouter(apiGroup)
}
