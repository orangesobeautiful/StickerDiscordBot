package controllers

import (
	"backend/pkg/hserr"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func (c *Controller) WebUserAuthRequired(ctx *gin.Context) {
	userID := sessions.Default(ctx).Get("user-auth")
	if userID == nil {
		ctx.AbortWithStatusJSON(401, hserr.ErrUnauthorized)
		return
	}
}
