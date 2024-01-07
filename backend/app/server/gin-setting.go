package server

import (
	"reflect"

	"backend/app/config"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"golang.org/x/text/language"
)

func newGinEngine(corsCfg config.CORS, validate *validator.Validate, eh *errHandler) *gin.Engine {
	setGinGlobal(validate)
	setGinextErrorHandler(eh)

	e := gin.Default()
	setGinLangDeal(e)
	setGinCORS(e, corsCfg)

	return e
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

func setGinCORS(e *gin.Engine, corsCfg config.CORS) {
	if corsCfg != nil {
		var allowOriginFunc func(origin string) bool
		if len(corsCfg.GetAllowOrigins()) > 0 {
			allowOriginSet := mapset.NewSet[string]()
			for _, origin := range corsCfg.GetAllowOrigins() {
				allowOriginSet.Add(origin)
			}
			allowOriginFunc = func(origin string) bool {
				return allowOriginSet.Contains(origin)
			}
		}

		e.Use(
			cors.New(cors.Config{
				AllowOriginFunc:  allowOriginFunc,
				AllowMethods:     corsCfg.GetAllowMethods(),
				AllowHeaders:     corsCfg.GetAllowHeaders(),
				ExposeHeaders:    corsCfg.GetExposeHeaders(),
				AllowCredentials: corsCfg.GetAllowCredentials(),
				MaxAge:           corsCfg.GetMaxAge(),
			}))
	}
}
