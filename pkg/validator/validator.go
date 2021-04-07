package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator"
)

type Validator struct {
	v *validator.Validate
}

// CreateValidator : Returns a new validator that can extract json tag values from struct
func CreateValidator() *Validator {
	val := new(Validator)
	val.v = validator.New()

	val.v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	return val
}

//	ValidateRequired : Returns a list of empty struct fields that have the required tag
//	The returned list is a list a strings of the struct fields json tags
func (val *Validator) ValidateRequired(s interface{}) []string {
	missingFields := make([]string, 0)

	valErrList := val.v.Struct(s)
	if valErrList != nil {
		for _, e := range valErrList.(validator.ValidationErrors) {
			if e.Tag() == "required" {
				missingFields = append(missingFields, e.Field())
			}
		}
	}

	return missingFields
}
