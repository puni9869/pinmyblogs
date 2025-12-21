package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/puni9869/pinmyblogs/pkg/formbinding"
)

const (
	contextKey = "pinmyblogs"
	formKey    = "__form"
)

func GetForm(c *gin.Context) any {
	if f, ok := c.Get(formKey); ok {
		return f
	}
	return nil
}

func GetContext(c *gin.Context) gin.H {
	if ctx, ok := c.Get(contextKey); ok {
		return ctx.(map[string]any)
	}
	return nil
}

// Bind binding a form obj to a handler's context data
func Bind[T any](_ T) gin.HandlerFunc {
	return func(c *gin.Context) {
		data := make(map[string]any)
		data["HasError"] = false
		var theObj = new(T) // create a new form obj for every request but not use obj directly
		c.Set(contextKey, data)
		var errs error
		if c.Request.Header.Get("Content-Type") == binding.MIMEJSON {
			errs = c.ShouldBindWith(theObj, binding.JSON)
		} else {
			errs = c.ShouldBindWith(theObj, binding.Form)
		}

		formbinding.FillContext(theObj, data)
		if errs != nil {
			//nolint:errorlint // errs is guaranteed to be validator.ValidationErrors here
			data = formbinding.Errorf(make(gin.H), errs.(validator.ValidationErrors))
			c.Set(contextKey, data)
		}
		c.Set(formKey, theObj)
	}
}
