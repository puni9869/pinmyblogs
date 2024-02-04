package formbinding

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type FieldErrors map[string]any

var fieldErrors = map[string]string{
	"required":           ` cannot be empty.`,
	"alpha_dash":         ` should contain only alphanumeric, dash ('-') and underscore ('_') characters.`,
	"alpha_dash_dot":     ` should contain only alphanumeric, dash ('-'), underscore ('_') and dot ('.') characters.`,
	"size":               ` must be size %s.`,
	"min_size":           ` must contain at least %s characters.`,
	"max_size":           ` must contain at most %s characters.`,
	"email":              ` is not a valid email address.`,
	"url":                `"%s" is not a valid URL.`,
	"include":            ` must contain substring "%s".`,
	"password_not_match": ` password not matched.`,
	"username":           ` can only contain alphanumeric chars ('0-9','a-z','A-Z'), dash ('-'), underscore ('_') and dot ('.'). It cannot begin or end with non-alphanumeric chars, and consecutive non-alphanumeric chars are also forbidden.`,
	"unknown":            `Unknown error:`,
}

func (e FieldErrors) Error(phrase string) string {
	if e, ok := fieldErrors[phrase]; ok {
		return e
	}
	return phrase
}

// Errorf format the error in the gin.H which is a context.
func Errorf(data gin.H, errs validator.ValidationErrors) map[string]any {
	if len(errs) == 0 {
		return nil
	}
	data["HasError"] = true
	var f FieldErrors
	for _, err := range errs {
		fmt.Println(err.Field())
		data[err.Field()+"_Error"] = fmt.Sprintf("%s%s", err.Field(), f.Error(err.Tag()))
		data[err.Field()+"_HasError"] = true
	}
	return data
}
