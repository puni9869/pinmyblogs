package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/puni9869/pinmyblogs/internal/signup"
	"github.com/puni9869/pinmyblogs/pkg/formbinding"
	"github.com/puni9869/pinmyblogs/server/middlewares"
	"github.com/puni9869/pinmyblogs/types/forms"
	"net/http"
)

var (
	fieldError formbinding.FieldErrors
)

func SignupGet(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.tmpl", nil)
}

func SignupPost(s signup.SignupService) gin.HandlerFunc {
	return func(c *gin.Context) {
		signup := middlewares.GetForm(c).(*forms.SignUpForm)
		ctx := middlewares.GetContext(c)

		password := signup.Password
		confirmPassword := signup.ConfirmPassword

		if ctx["Password_HasError"] == false && len(password) != len(confirmPassword) {
			ctx["Password_Error"] = ""
			ctx["Password_HasError"] = true
			ctx["ConfirmPassword_Error"] = fieldError.Error("password_not_match")
			ctx["ConfirmPassword_HasError"] = true
		}

		c.HTML(http.StatusOK, "signup.tmpl", ctx)
	}
}
