// Used for @error validation
// See views/components/form/error.templ

package utils

import (
	"context"
	"reflect"

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
