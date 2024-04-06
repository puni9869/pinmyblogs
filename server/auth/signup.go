package auth

import (
	"crypto/sha256"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/puni9869/pinmyblogs/internal/signup"
	"github.com/puni9869/pinmyblogs/models"
	"github.com/puni9869/pinmyblogs/pkg/formbinding"
	"github.com/puni9869/pinmyblogs/pkg/logger"
	"github.com/puni9869/pinmyblogs/server/middlewares"
	"github.com/puni9869/pinmyblogs/types/forms"
	"net/http"
	"sync"
)

var (
	field formbinding.Field
	once  sync.Once
)

// SignupGet is renders the signup.tmpl
// c is gin.Context
func SignupGet(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.tmpl", nil)
}

// SignupPost is a handler for handling registrations of new user
// signUp is signup.Service container helper functions
func SignupPost(signUp signup.Service) gin.HandlerFunc {
	log := logger.NewLogger()
	return func(c *gin.Context) {
		form := middlewares.GetForm(c).(*forms.SignUpForm)
		ctx := middlewares.GetContext(c)

		// initialize only once
		once.Do(func() {
			field = new(formbinding.FieldErrors)
		})
		password := form.Password
		email := form.Email
		confirmPassword := form.ConfirmPassword

		// password check
		if ctx["Email_HasError"] == false && !field.IsValid(password) {
			ctx["HasError"] = true
			ctx["Password_HasError"] = true
			ctx["Password_Error"] = field.Error("alpha_dash_dot")
			ctx["ConfirmPassword_HasError"] = false
		}

		// password and confirm password checks
		if (ctx["Password_HasError"] == nil || ctx["Password_HasError"] == false) &&
			(len(password) != len(confirmPassword) || password != confirmPassword) {
			ctx["HasError"] = true
			ctx["Password_Error"] = ""
			ctx["Password_HasError"] = true
			ctx["ConfirmPassword_Error"] = field.Error("password_not_match")
			ctx["ConfirmPassword_HasError"] = true
		}

		if ctx["HasError"] == false {
			// using sha256 hash getting the checksums i.e. one way hash for password
			h := sha256.New()
			h.Write([]byte(password))
			bs := h.Sum(nil)

			user := models.User{Password: fmt.Sprintf("%x", bs), Email: email}
			err := signUp.Register(c, user)
			if err == nil {
				log.WithFields(map[string]any{
					"user":      user.Email,
					"id":        user.ID,
					"createdAt": user.CreatedAt,
				}).Info("user is registered")
			}
			ctx["HasError"] = true
			log.WithError(err).Error("error in registering user")
		}
		c.HTML(http.StatusOK, "signup.tmpl", ctx)
	}
}
