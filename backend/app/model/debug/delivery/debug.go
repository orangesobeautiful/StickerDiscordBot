package delivery

import (
	ginauth "backend/app/model/discorduser/gin-auth"

	"github.com/gin-gonic/gin"
)

type debugController struct {
	auth ginauth.AuthInterface
}

func Initialze(
	e *gin.Engine,
	auth ginauth.AuthInterface,
) {
	ctrl := debugController{
		auth: auth,
	}

	ctrl.RegisterGinRouter(e)
}
