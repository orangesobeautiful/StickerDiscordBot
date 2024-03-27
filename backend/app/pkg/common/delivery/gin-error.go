package delivery

import "github.com/gin-gonic/gin"

func GinAbortWithError(ctx *gin.Context, err error) {
	_ = ctx.Error(err)
	ctx.Abort()
}
