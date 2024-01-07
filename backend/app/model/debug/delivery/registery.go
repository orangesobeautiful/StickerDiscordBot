package delivery

import (
	"github.com/arl/statsviz"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func (c *debugController) RegisterGinRouter(e *gin.Engine) {
	authGroup := e.Group("")
	authGroup.Use(c.auth.GetRequiredAuthMiddleware())

	debugGroup := e.Group("/debug")

	pprof.RouteRegister(debugGroup, "pprof")

	c.registerStatsvizGinRouter(debugGroup)
}

func (c *debugController) registerStatsvizGinRouter(debugGroup *gin.RouterGroup) {
	statsvizGroup := debugGroup.Group("/statsviz")

	srv, err := statsviz.NewServer(statsviz.Root("/debug/statsviz"))
	if err != nil {
		panic(err)
	}

	index := srv.Index()
	ws := srv.Ws()

	statsvizGroup.GET("/*filepath", func(context *gin.Context) {
		if context.Param("filepath") == "/ws" {
			ws(context.Writer, context.Request)
			return
		}
		index(context.Writer, context.Request)
	})
}
