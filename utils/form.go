// Used for @error validation
// See views/components/form/error.templ

package utils

import (
	"context"
	"reflect"
	"unicode"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

func WithErrors(ctx context.Context, errs map[string]string) context.Context {
	return context.WithValue(ctx, "errors", errs)
}

func ErrorsFrom(ctx context.Context) map[string]string {
	if v := ctx.Value("errors"); v != nil {
		if errs, ok := v.(map[string]string); ok {
			return errs
		}
	}
	return map[string]string{}
}

func ValidateStruct(s any) (map[string]string, bool) {
	errs := map[string]string{}
	validate := validator.New()
	uni := ut.New(en.New())
	ENTranslation, _ := uni.GetTranslator("en")
	_ = en_translations.RegisterDefaultTranslations(validate, ENTranslation)
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("form")
		if name == "" {
			return fld.Name
		}
		return name
	})
	validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		var hasUpper, hasLower, hasDigit, hasSpecial bool

		for _, ch := range password {
			switch {
			case unicode.IsUpper(ch):
				hasUpper = true
			case unicode.IsLower(ch):
				hasLower = true
			case unicode.IsDigit(ch):
				hasDigit = true
			case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
				hasSpecial = true
			}
		}

		return hasUpper && hasLower && hasDigit && hasSpecial
	})
	validate.RegisterTranslation("password", ENTranslation, func(ut ut.Translator) error {
		return ut.Add("password", "{0} must contain at least 1 uppercase, 1 lowercase, 1 number, and 1 special character", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("password", fe.Field())
		return t
	})
	err := validate.Struct(s)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			field := e.Field()
			//errs[field] = e.Error()
			errs[field] = e.Translate(ENTranslation)
		}
	}
	return errs, len(errs) == 0
}
