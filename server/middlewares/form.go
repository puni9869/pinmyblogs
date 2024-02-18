package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/puni9869/pinmyblogs/pkg/formbinding"
)

const (
	ContextKey = "pinmyblogs"
	formKey    = "__form"
)

func GetForm(c *gin.Context) any {
	if f, ok := c.Get(formKey); ok {
		return f
	}
	return nil
}

func GetContext(c *gin.Context) gin.H {
	if ctx, ok := c.Get(ContextKey); ok {
		return ctx.(map[string]any)
	}
	return nil
}

// BindForm binding a form obj to a handler's context data
func BindForm[T any](_ T) gin.HandlerFunc {
	return func(c *gin.Context) {
		data := make(map[string]any)
		data["HasError"] = false
		c.Set(ContextKey, data)
		var theObj = new(T) // create a new form obj for every request but not use obj directly
		if errs := c.ShouldBindWith(theObj, binding.Form); errs != nil {
			data = formbinding.Errorf(make(gin.H), errs.(validator.ValidationErrors))
			c.Set(ContextKey, data)
		}
		c.Set(formKey, theObj)
	}
}
