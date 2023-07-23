package server

import (
	"fmt"
	"reflect"

	"backend/config"
	"backend/pkg/log"
	"backend/server/controllers"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en_US"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"golang.org/x/text/language"
)

type Server struct {
	*gin.Engine

	cfg           *config.CfgInfo
	uniTranslator *ut.UniversalTranslator
	langMatcher   language.Matcher

	ctrl *controllers.Controller
}

func New(cfg *config.CfgInfo) (*Server, error) {
	var err error

	// init gin engine

	e, err := newGinEngine(cfg)
	if err != nil {
		log.Errorf("newEngine failed, err=%s", err)
		return nil, err
	}

	return e, nil
}

func newGinEngine(cfg *config.CfgInfo) (*Server, error) {
	e := gin.Default()

	var err error
	var uni *ut.UniversalTranslator
	if uni, err = setGinValidate(); err != nil {
		return nil, fmt.Errorf("setValidate failed: %w", err)
	}

	s := &Server{
		Engine:        e,
		cfg:           cfg,
		uniTranslator: uni,
		ctrl:          controllers.New(cfg),
	}
	s.setGinRouter()

	return s, nil
}

func setGinValidate() (*ut.UniversalTranslator, error) {
	validate := binding.Validator.Engine().(*validator.Validate)

	english := en_US.New()
	uni := ut.New(english, english)
	en, _ := uni.GetTranslator("en")

	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		transTag := field.Tag.Get("trans")
		if transTag == "" {
			return field.Name
		}
		return transTag
	})

	err := en_translations.RegisterDefaultTranslations(validate, en)
	if err != nil {
		return nil, err
	}

	return uni, nil
}
