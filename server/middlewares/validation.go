package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/puni9869/pinmyblogs/pkg/formbinding"
)

const (
	formErrorKey = "FormErrorKey"
	formKey      = "__form"
)

func GetForm(c *gin.Context) any {
	if f, ok := c.Get(formKey); ok {
		return f
	}
	return nil
}

func GetFormErr(c *gin.Context) any {
	if errs, ok := c.Get(formErrorKey); ok {
		return errs
	}
	return nil
}

// BindForm binding a form obj to a handler's context data
func BindForm[T any](_ T) gin.HandlerFunc {
	return func(c *gin.Context) {
		var theObj T // create a new form obj for every request but not use obj directly
		if errs := c.ShouldBindWith(&theObj, binding.Form); errs != nil {
			eData := formbinding.Errorf(errs.(validator.ValidationErrors))
			c.Set(formErrorKey, eData)
		}
		c.Set(formKey, theObj)
	}
}
