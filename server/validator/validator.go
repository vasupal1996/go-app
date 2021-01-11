package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/validator/v10"

	ut "github.com/go-playground/universal-translator"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	errors "github.com/vasupal1996/goerror"
)

// Validator container validator library and transaltor
type Validator struct {
	V *validator.Validate
	T *ut.Translator
}

// NewValidation create new Validator struct instance
func NewValidation() *Validator {
	v := &Validator{
		V: validator.New(),
	}
	trans := initializeTranslation(v.V)
	v.T = trans
	registerFunc(v.V)
	return v
}

// Initialize initializes and returns the UniversalTranslator instance for the application
func initializeTranslation(validate *validator.Validate) *ut.Translator {

	// initialize translator
	en := en.New()
	uni := ut.New(en, en)

	trans, _ := uni.GetTranslator("en")
	// initialize translations
	en_translations.RegisterDefaultTranslations(validate, trans)
	return &trans
}

func registerFunc(validate *validator.Validate) {
	// register function to get tag name from json tags.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// Validate validates the struct
// Note: do not pass slice of struct
func (v *Validator) Validate(form interface{}) []error {
	var errResp []error
	if err := v.V.Struct(form); err != nil {
		errs := err.(validator.ValidationErrors)
		for _, e := range errs {
			err := errors.New(e.Translate(*v.T), &errors.BadRequest)
			key := strings.SplitAfterN(e.Namespace(), ".", 2)
			err = errors.SetContext(err, key[1], e.Translate(*v.T))
			errResp = append(errResp, err)
		}
	}
	return errResp
}
