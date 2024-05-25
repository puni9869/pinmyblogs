package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/puni9869/pinmyblogs/models"
	"github.com/puni9869/pinmyblogs/pkg/config"
	"github.com/puni9869/pinmyblogs/pkg/database"
	"github.com/puni9869/pinmyblogs/pkg/logger"
	"github.com/puni9869/pinmyblogs/pkg/mailer"
	"github.com/puni9869/pinmyblogs/server/middlewares"
	"github.com/puni9869/pinmyblogs/types/forms"
)

func ResetPasswordGet(c *gin.Context) {
	c.HTML(http.StatusOK, "reset.tmpl", nil)
}

func ResetPasswordPost(c *gin.Context) {
	log := logger.NewLogger()

	if !config.C.Authentication.EnableForgotPassword {
		log.
			WithField("isEnableForgotPassword", config.C.Authentication.EnableForgotPassword).
			Warn("Forgot Password is disabled globally.")

		c.HTML(http.StatusOK, "reset.tmpl", gin.H{
			"HasError": true,
			"Error":    "Reset password is currently disable.",
		})
		c.Abort()
		return
	}
	form := middlewares.GetForm(c).(*forms.ResetForm)
	ctx := middlewares.GetContext(c)

	email := form.Email

	if ctx["Email_HasError"] == true {
		fmt.Print("Here")
		log.WithField("email", email).Error("email id not found.")
		c.HTML(http.StatusBadRequest, "reset.tmpl", gin.H{"Email": email, "HasError": true, "Error": "Email id not found."})
		return
	}

	var user *models.User
	result := database.Db().First(&user, "email = ?", email)
	if result.Error != nil {
		log.WithField("email", email).WithError(result.Error).Error("Invalid email or password. Database error")
		c.HTML(http.StatusBadRequest, "reset.tmpl", gin.H{"Email": email, "HasError": true, "Error": "Email id not found."})
		c.Abort()
		return
	}

	if !user.IsActive || !user.IsEmailVerified {
		log.WithFields(map[string]any{
			"email":           user.Email,
			"isActive":        user.IsActive,
			"isEmailVerified": user.IsEmailVerified,
		}).WithError(result.Error).Error("Account is disabled.")
		c.HTML(http.StatusUnauthorized, "reset.tmpl", gin.H{"HasError": true, "Error": "Account is disabled."})
		c.Abort()
		return
	}
	log.WithField("email", user.Email).Info("Email account found for password reset.")

	resetMailer := mailer.NewResetPasswordMailer(*user)
	go resetMailer.Send()

	c.HTML(http.StatusAccepted, "reset_password_sent.tmpl", gin.H{"Email": email})
}
