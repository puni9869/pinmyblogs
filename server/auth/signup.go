package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/puni9869/pinmyblogs/internal/signup"
	"github.com/puni9869/pinmyblogs/server/middlewares"
	"github.com/puni9869/pinmyblogs/types/forms"
	"net/http"
)

func SignupGet(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.tmpl", nil)
}

func SignupPost(signupService signup.SignupService) gin.HandlerFunc {
	return func(c *gin.Context) {
		_ = middlewares.GetForm(c).(forms.SignUpForm)
		errs := middlewares.GetFormErr(c)
		c.HTML(http.StatusOK, "signup.tmpl", gin.H{"errors": errs})
	}
}
