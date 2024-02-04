package formbinding

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

var fieldErrors = map[string]string{
	"required":       ` cannot be empty.`,
	"alpha_dash":     ` should contain only alphanumeric, dash ('-') and underscore ('_') characters.`,
	"alpha_dash_dot": ` should contain only alphanumeric, dash ('-'), underscore ('_') and dot ('.') characters.`,
	"size":           ` must be size %s.`,
	"min_size":       ` must contain at least %s characters.`,
	"max_size":       ` must contain at most %s characters.`,
	"email":          ` is not a valid email address.`,
	"url":            `"%s" is not a valid URL.`,
	"include":        ` must contain substring "%s".`,
	"username":       ` can only contain alphanumeric chars ('0-9','a-z','A-Z'), dash ('-'), underscore ('_') and dot ('.'). It cannot begin or end with non-alphanumeric chars, and consecutive non-alphanumeric chars are also forbidden.`,
	"unknown":        `Unknown error:`,
}

// Errorf validates the form
func Errorf(errs validator.ValidationErrors) map[string]any {
	if len(errs) == 0 {
		return nil
	}
	var data = make(map[string]any)
	data["HasError"] = true
	for _, err := range errs {
		if phrase, ok := fieldErrors[err.Tag()]; ok {
			data[err.Field()+"_Err"] = fmt.Sprintf("%s%s", err.Field(), phrase)
			data[err.Field()+"_HasError"] = true
		}
	}
	return data
}
