package server

import (
	"encoding/json"
	"net/http"

	"backend/pkg/ginext"
	"backend/pkg/hserr"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/text/language"
)

func (s *Server) setGinRouter() {
	e := s.Engine

	s.setCORS()
	s.setLangDeal()
	s.setErrorHandler()
	sessStore := s.newSessStore()

	// Serve frontend static files (SPA mode)
	e.Use(static.Serve("/", static.LocalFile("public", false)))
	e.NoRoute(func(ctx *gin.Context) {
		ctx.File("public/index.html")
	})

	// Serve sticker image
	e.Static("/sticker-image/", "sticker-image")

	// Web API
	webAPIGroup := e.Group("/api/web")
	webAPIV1Group := webAPIGroup.Group("/v1")

	publicGroup := webAPIV1Group
	publicGroup.GET("/gen_login_code", ginext.Handler(s.ctrl.WebGenLoginCode))
	publicGroup.GET("/check_login",
		sessions.Sessions("user-auth", sessStore),
		ginext.BindHandler(s.ctrl.WebCheckLogin))

	authRequiredRouter := webAPIV1Group.Use(
		sessions.Sessions("user-auth", sessStore), s.ctrl.WebUserAuthRequired)
	authRequiredRouter.GET("/logout", ginext.Handler(s.ctrl.WebLogout))
	authRequiredRouter.GET("/has_login", ginext.Handler(s.ctrl.WebHasLlogin))
	authRequiredRouter.GET("/user_info", ginext.Handler(s.ctrl.WebSelfInfo))

	authRequiredRouter.GET("/all_sticker", ginext.BindHandler(s.ctrl.ListSticker))
	authRequiredRouter.GET("/search", ginext.BindHandler(s.ctrl.SearchSticker))
	authRequiredRouter.POST("/change_sn", ginext.BindHandler(s.ctrl.ChangeSticker))
}

func (s *Server) newSessStore() sessions.Store {
	cookieStore := cookie.NewStore(s.cfg.Server.SessionKey.UserAuth.SessionKeyPair()...)
	if s.cfg.Server.Cookie != nil {
		cookieStore.Options(sessions.Options{
			MaxAge:   int(s.cfg.Server.CORS.MaxAge.Seconds()),
			Secure:   s.cfg.Server.Cookie.Secure,
			HttpOnly: s.cfg.Server.Cookie.HttpOnly,
			SameSite: s.cfg.Server.Cookie.SameSite,
		})
	}

	return cookieStore
}

func (s *Server) setCORS() {
	if s.cfg.Server.CORS != nil {
		s.Engine.Use(
			cors.New(cors.Config{
				AllowOrigins:     s.cfg.Server.CORS.AllowOrigins,
				AllowMethods:     s.cfg.Server.CORS.AllowMethods,
				AllowHeaders:     s.cfg.Server.CORS.AllowHeaders,
				ExposeHeaders:    s.cfg.Server.CORS.ExposeHeaders,
				AllowCredentials: s.cfg.Server.CORS.AllowCredentials,
				MaxAge:           s.cfg.Server.CORS.MaxAge,
			}))
	}
}

func (s *Server) setLangDeal() {
	s.langMatcher = language.NewMatcher([]language.Tag{
		language.English,
	})

	s.Engine.Use(func(ctx *gin.Context) {
		// TODO: layz init langTag

		lang, _ := ctx.Cookie("lang")
		accept := ctx.Request.Header.Get("Accept-Language")
		tag, _ := language.MatchStrings(s.langMatcher, lang, accept)

		ctx.Set("lang", tag)
	})
}

func (s *Server) setErrorHandler() {
	ginext.SetBindErrHandler(func(ctx *gin.Context, err error) {
		switch realErr := err.(type) {
		case validator.ValidationErrors:
			trans, _ := s.uniTranslator.GetTranslator(
				ctx.MustGet("lang").(language.Tag).String())

			detailList := make([]string, 0, len(realErr))
			for _, valErr := range realErr {
				detailList = append(detailList, valErr.Translate(trans))
			}

			ctx.JSON(http.StatusBadRequest, hserr.ErrResp{
				Message: "param of request validate failed",
				Detail:  detailList,
			})
			return
		case *json.UnmarshalTypeError:
			ctx.JSON(http.StatusBadRequest, hserr.ErrResp{
				Message: "decode json body failed",
				Detail: []string{
					realErr.Field +
						" should be " +
						realErr.Type.Name() +
						" not " + realErr.Value,
				},
			})
		case *json.SyntaxError:

			ctx.JSON(http.StatusBadRequest, hserr.ErrResp{
				Message: "decode json body failed",
				Detail:  []string{err.Error()},
			})
		default:
			ctx.JSON(http.StatusBadRequest, hserr.ErrResp{
				Message: "bad request format",
				Detail:  []string{err.Error()},
			})
		}
	})

	ginext.SetRespErrHandler(func(ctx *gin.Context, err error) {
		switch realErr := err.(type) {
		case *hserr.ErrResp:
			ctx.JSON(realErr.HttpStatus, realErr)
			return
		default:
			ctx.JSON(http.StatusInternalServerError, hserr.ErrResp{
				Message: "unknown error",
				Detail:  []string{err.Error()},
			})
			return
		}
	})
}
