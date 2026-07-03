package models

import (
	"reflect"
	"strings"

	en "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var (
	validate *validator.Validate
	trans    ut.Translator
)

func init() {
	validate = validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("label"), ",", 2)[0]
		if name != "" {
			return name
		}
		return fld.Name
	})

	english := en.New()
	uni := ut.New(english, english)
	trans, _ = uni.GetTranslator("en")
	_ = en_translations.RegisterDefaultTranslations(validate, trans)
}

func (input CreateUserInput) Validate() error {
	return validate.Struct(input)
}

func FormatValidationError(err error) string {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err.Error()
	}

	messages := make([]string, 0, len(validationErrors))
	for _, fieldErr := range validationErrors {
		messages = append(messages, fieldErr.Translate(trans))
	}
	return strings.Join(messages, "; ")
}
