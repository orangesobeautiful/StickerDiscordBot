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
	"golang.org/x/text/language"
)

func newGinEngine(cfg *config.CfgInfo, validate *validator.Validate, uni *ut.UniversalTranslator) *gin.Engine {
	setGinGlobal(validate)
	setGinextErrorHandler(uni)

	e := gin.Default()
	setGinLangDeal(e)
	setGinCORS(e, cfg)

	return e
}

func newValidateTranslator() *ut.UniversalTranslator {
	english := en_US.New()
	uni := ut.New(english, english)

	return uni
}

func setGinGlobal(validate *validator.Validate) {
	setGinValidate(validate)
}

type ginBindingValidator struct {
	validate *validator.Validate
}

// ValidateStruct receives any kind of type, but only performed struct or pointer to struct type.
func (v *ginBindingValidator) ValidateStruct(obj any) error {
	if obj == nil {
		return nil
	}

	value := reflect.ValueOf(obj)
	switch value.Kind() {
	case reflect.Ptr:
		return v.ValidateStruct(value.Elem().Interface())
	case reflect.Struct:
		return v.validateStruct(obj)
	case reflect.Slice, reflect.Array:
		count := value.Len()
		validateRet := make(binding.SliceValidationError, 0)
		for i := 0; i < count; i++ {
			if err := v.ValidateStruct(value.Index(i).Interface()); err != nil {
				validateRet = append(validateRet, err)
			}
		}
		if len(validateRet) == 0 {
			return nil
		}
		return validateRet
	default:
		return nil
	}
}

// validateStruct receives struct type
func (v *ginBindingValidator) validateStruct(obj any) error {
	return v.validate.Struct(obj)
}

func (v *ginBindingValidator) Engine() any {
	return v.validate
}

func setGinValidate(validate *validator.Validate) {
	binding.Validator = &ginBindingValidator{validate: validate}
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
