package hs

import (
	"reflect"

	"github.com/go-playground/locales/en_US"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var uni *ut.UniversalTranslator

func NewDefaultValidate() (*validator.Validate, error) {
	english := en_US.New()

	uni = ut.New(english, english)

	en, _ := uni.GetTranslator("en")

	var err error
	validate := validator.New()
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		transTag := field.Tag.Get("trans")
		if transTag == "" {
			return field.Name
		}
		return transTag
	})

	err = en_translations.RegisterDefaultTranslations(validate, en)
	if err != nil {
		return nil, err
	}

	return validate, nil
}
