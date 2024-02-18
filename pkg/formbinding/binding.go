package formbinding

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"reflect"
	"regexp"
)

const (
	regexPattern = `^[A-Za-z\d\W_]{6,15}$`
)

type Field interface {
	Error(phrase string) string
	IsValid(password string) bool
}

type FieldErrors map[string]string

var errs = map[string]string{
	"required":           `Cannot be empty`,
	"alpha_dash":         `Must contain only alphanumeric, dash ('-') and underscore ('_') characters.`,
	"alpha_dash_dot":     `Mandate 1 letter, 1 digit, 7-15 characters, with 6 non-digits`,
	"size":               `must be size %s.`,
	"min_size":           `must contain at least %s characters.`,
	"max_size":           `must contain at most %s characters.`,
	"email":              `is not a valid email address.`,
	"url":                `"%s" is not a valid URL.`,
	"include":            `must contain substring "%s".`,
	"password_not_match": `Password not match`,
	"username":           `can only contain alphanumeric chars ('0-9','a-z','A-Z'), dash ('-'), underscore ('_') and dot ('.'). It cannot begin or end with non-alphanumeric chars, and consecutive non-alphanumeric chars are also forbidden.`,
	"unknown":            `Unknown error:`,
}

func (e *FieldErrors) Error(phrase string) string {
	if e, ok := errs[phrase]; ok {
		return e
	}
	return phrase
}

func (e *FieldErrors) IsValid(password string) bool {
	match, err := regexp.MatchString(regexPattern, password)
	if err != nil {
		panic(err)
	}
	return match
}

// Errorf format the error in the gin.H which is a context.
func Errorf(data gin.H, errs validator.ValidationErrors) map[string]any {
	if len(errs) == 0 {
		return nil
	}
	var f Field = new(FieldErrors)
	for _, err := range errs {
		data[err.Field()+"_Error"] = fmt.Sprintf("%s%s", err.Field(), f.Error(err.Tag()))
		data[err.Field()+"_HasError"] = true
	}
	return data
}

// FillContext is fill all the form fields into *gin.Context
//
// data[field.Name+"_Error"] = ""
// data[field.Name+"_HasError"] = false
func FillContext(form any, data gin.H) {
	typ := reflect.TypeOf(form)
	val := reflect.ValueOf(form)
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		val = val.Elem()
	}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldName := field.Tag.Get("form")
		// Allow ignored fields in the struct
		if fieldName == "-" {
			continue
		} else if len(fieldName) == 0 {
			fieldName = field.Name
		}
		data[fieldName] = val.Field(i).Interface()
		data[field.Name+"_Error"] = ""
		data[field.Name+"_HasError"] = false
	}
}
