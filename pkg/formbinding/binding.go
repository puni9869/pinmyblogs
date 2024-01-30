package formbinding

import (
	"fmt"
	"reflect"
	"strings"

	"gitea.com/go-chi/binding"
	"github.com/puni9869/pinmyblogs/pkg/util"
)

const (
	// ErrUsername is username error
	ErrUsername = "UsernameError"
)

func init() {
	binding.SetNameMapper(util.ToSnakeCase)
}

// Form form binding interface
type Form interface {
	binding.Validator
}

// AssignForm assign form values back to the template data.
func AssignForm(form any, data map[string]any) {
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
			fieldName = util.ToSnakeCase(field.Name)
		}

		data[fieldName] = val.Field(i).Interface()
	}
}

func getRuleBody(field reflect.StructField, prefix string) string {
	for _, rule := range strings.Split(field.Tag.Get("binding"), ";") {
		if strings.HasPrefix(rule, prefix) {
			return rule[len(prefix) : len(rule)-1]
		}
	}
	return ""
}

// GetSize get size int form tag
func GetSize(field reflect.StructField) string {
	return getRuleBody(field, "Size(")
}

// GetMinSize get minimal size in form tag
func GetMinSize(field reflect.StructField) string {
	return getRuleBody(field, "MinSize(")
}

// GetMaxSize get max size in form tag
func GetMaxSize(field reflect.StructField) string {
	return getRuleBody(field, "MaxSize(")
}

// GetInclude get include in form tag
func GetInclude(field reflect.StructField) string {
	return getRuleBody(field, "Include(")
}

// Validate validates the form
func Validate(errs binding.Errors, f Form) binding.Errors {
	fmt.Printf("Hrre")
	var data = make(map[string]any)
	if errs.Len() == 0 {
		return errs
	}

	data["HasError"] = true
	// If the field with name errs[0].FieldNames[0] is not found in form
	// somehow, some code later on will panic on Data["ErrorMsg"].(string).
	// So initialize it to some default.
	data["ErrorMsg"] = unknown_error
	AssignForm(f, data)

	typ := reflect.TypeOf(f)

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if field, ok := typ.FieldByName(errs[0].FieldNames[0]); ok {
		fieldName := field.Tag.Get("form")
		if fieldName != "-" {
			data["Err_"+field.Name] = true

			trName := field.Tag.Get("locale")
			if len(trName) == 0 {
				trName = "form." + field.Name
			} else {
				trName = trName
			}

			switch errs[0].Classification {
			case binding.ERR_REQUIRED:
				data["ErrorMsg"] = trName + require_error
			case binding.ERR_ALPHA_DASH:
				data["ErrorMsg"] = trName + alpha_dash_error
			case binding.ERR_ALPHA_DASH_DOT:
				data["ErrorMsg"] = trName + alpha_dash_dot_error
			case binding.ERR_SIZE:
				data["ErrorMsg"] = trName + formatErr(size_error, GetSize(field))
			case binding.ERR_MIN_SIZE:
				data["ErrorMsg"] = trName + formatErr(min_size_error, GetMinSize(field))
			case binding.ERR_MAX_SIZE:
				data["ErrorMsg"] = trName + formatErr(max_size_error, GetMaxSize(field))
			case binding.ERR_EMAIL:
				data["ErrorMsg"] = trName + email_error
			case binding.ERR_URL:
				data["ErrorMsg"] = trName + formatErr(url_error, errs[0].Message)
			case binding.ERR_INCLUDE:
				data["ErrorMsg"] = trName + formatErr(include_error, GetInclude(field))
			case ErrUsername:
				data["ErrorMsg"] = trName + username_error

			default:
				msg := errs[0].Classification
				if msg != "" && errs[0].Message != "" {
					msg += ": "
				}

				msg += errs[0].Message
				if msg == "" {
					msg = unknown_error
				}
				data["ErrorMsg"] = trName + ": " + msg
			}
			return errs
		}
	}
	return errs
}
