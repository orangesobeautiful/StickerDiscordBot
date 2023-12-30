package server

import (
	"reflect"

	"backend/app/config"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en_US"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"golang.org/x/text/language"
	"golang.org/x/xerrors"
)

func newGinEngine(cfg *config.CfgInfo) (*gin.Engine, error) {
	var err error

	uni := newValidateTranslator()
	err = setGinGlobal(uni)
	if err != nil {
		return nil, err
	}
	setGinextErrorHandler(uni)

	e := gin.Default()
	setGinLangDeal(e)
	setGinCORS(e, cfg)

	return e, nil
}

func newValidateTranslator() *ut.UniversalTranslator {
	english := en_US.New()
	uni := ut.New(english, english)

	return uni
}

func setGinGlobal(uni *ut.UniversalTranslator) (err error) {
	err = setGinValidate(uni)
	if err != nil {
		return xerrors.Errorf("setGinValidate : %w", err)
	}

	return nil
}

func setGinValidate(uni *ut.UniversalTranslator) error {
	en, _ := uni.GetTranslator("en")

	validate := binding.Validator.Engine().(*validator.Validate)
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		transTag := field.Tag.Get("trans")
		if transTag == "" {
			return field.Name
		}
		return transTag
	})

	err := en_translations.RegisterDefaultTranslations(validate, en)
	if err != nil {
		return err
	}

	return nil
}

func setGinLangDeal(e *gin.Engine) {
	langMatcher := language.NewMatcher([]language.Tag{
		language.English,
	})

	e.Use(func(ctx *gin.Context) {
		// TODO: layz init langTag

		lang, _ := ctx.Cookie("lang")
		accept := ctx.Request.Header.Get("Accept-Language")
		tag, _ := language.MatchStrings(langMatcher, lang, accept)

		ctx.Set("lang", tag)
	})
}

func setGinCORS(e *gin.Engine, cfg *config.CfgInfo) {
	if cfg.Server.CORS != nil {
		var allowOriginFunc func(origin string) bool
		if len(cfg.Server.CORS.AllowOrigins) > 0 {
			allowOriginSet := mapset.NewSet[string]()
			for _, origin := range cfg.Server.CORS.AllowOrigins {
				allowOriginSet.Add(origin)
			}
			allowOriginFunc = func(origin string) bool {
				return allowOriginSet.Contains(origin)
			}
		}

		e.Use(
			cors.New(cors.Config{
				AllowOriginFunc:  allowOriginFunc,
				AllowMethods:     cfg.Server.CORS.AllowMethods,
				AllowHeaders:     cfg.Server.CORS.AllowHeaders,
				ExposeHeaders:    cfg.Server.CORS.ExposeHeaders,
				AllowCredentials: cfg.Server.CORS.AllowCredentials,
				MaxAge:           cfg.Server.CORS.MaxAge,
			}))
	}
}

func newSessStore(cfg *config.CfgInfo) sessions.Store {
	cookieStore := cookie.NewStore(cfg.Server.SessionKey.UserAuth.SessionKeyPair()...)
	if cfg.Server.Cookie != nil {
		cookieStore.Options(sessions.Options{
			MaxAge:   int(cfg.Server.CORS.MaxAge.Seconds()),
			Secure:   cfg.Server.Cookie.Secure,
			HttpOnly: cfg.Server.Cookie.HTTPOnly,
			SameSite: cfg.Server.Cookie.SameSite,
		})
	}

	return cookieStore
}
