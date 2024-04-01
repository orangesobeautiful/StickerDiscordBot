package delivery

import (
	"net/http"
	"strconv"

	"backend/app/pkg/hserr"

	"github.com/gin-gonic/gin"
)

func GetIDParamFromContext(ctx *gin.Context, param string) (int, error) {
	idStr := ctx.Param(param)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, hserr.New(http.StatusBadRequest, param+" is not a number")
	}

	return id, nil
}

func GinAbortWithError(ctx *gin.Context, err error) {
	_ = ctx.Error(err)
	ctx.Abort()
}
