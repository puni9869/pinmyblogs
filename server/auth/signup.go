package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/puni9869/pinmyblogs/internal/signup"
	"github.com/puni9869/pinmyblogs/pkg/formbinding"
	"github.com/puni9869/pinmyblogs/server/middlewares"
	"github.com/puni9869/pinmyblogs/types/forms"
	"net/http"
	"sync"
)

var (
	field formbinding.Field
	once  sync.Once
)

// SignupGet is renders the signup.templ
// c is gin.Context
func SignupGet(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.tmpl", nil)
}

// SignupPost is a handler for handling registrations of new user
// s is signup.SignupService container helper functions
func SignupPost(s signup.SignupService) gin.HandlerFunc {
	return func(c *gin.Context) {
		s := middlewares.GetForm(c).(*forms.SignUpForm)
		ctx := middlewares.GetContext(c)
		once.Do(func() {
			field = new(formbinding.FieldErrors)
		})
		password := s.Password
		confirmPassword := s.ConfirmPassword

		if field.ValidatePassword(password) {
			ctx["Password_HasError"] = true
			ctx["Password_Error"] = field.Error("alpha_dash_dot")
			ctx["ConfirmPassword_HasError"] = true
			ctx["ConfirmPassword_Error"] = field.Error("password_not_match")
		}

		if ctx["Password_HasError"] == false && len(password) != len(confirmPassword) {
			ctx["Password_Error"] = ""
			ctx["Password_HasError"] = true
			ctx["ConfirmPassword_Error"] = field.Error("password_not_match")
			ctx["ConfirmPassword_HasError"] = true
		}

		c.HTML(http.StatusOK, "signup.tmpl", ctx)
	}
}
