package response

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// Translate ...
func Translate(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("is %s!", e.Tag())
	case "email":
		return "Please input a valid email!"
	case "len":
		return fmt.Sprintf("should be exactly %s character(s) or %s item(s)!", e.Param(), e.Param())
	case "max":
		return fmt.Sprintf("may not be more than %s character(s) or %s item(s)!", e.Param(), e.Param())
	case "min":
		return fmt.Sprintf("should be at least %s character(s) or %s item(s)!", e.Param(), e.Param())
	case "alphanum":
		return "may just contains alphabet and numeric!"
	case "number":
		return "should be a number!"
	case "gt":
		return fmt.Sprintf("should be greater than %s!", e.Param())
	case "gte":
		return fmt.Sprintf("should be greater than or equal %s!", e.Param())
	case "lt":
		return fmt.Sprintf("should be less than %s!", e.Param())
	case "lte":
		return fmt.Sprintf("should be less than or equal %s!", e.Param())
	case "unique":
		return "has already been taken!"
	}
	return fmt.Sprint(e)
}
