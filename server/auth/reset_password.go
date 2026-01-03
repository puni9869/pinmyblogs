package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/puni9869/pinmyblogs/models"
	"github.com/puni9869/pinmyblogs/pkg/config"
	"github.com/puni9869/pinmyblogs/pkg/database"
	"github.com/puni9869/pinmyblogs/pkg/formbinding"
	"github.com/puni9869/pinmyblogs/pkg/logger"
	"github.com/puni9869/pinmyblogs/pkg/mailer"
	"github.com/puni9869/pinmyblogs/pkg/utils"
	"github.com/puni9869/pinmyblogs/server/middlewares"
	"github.com/puni9869/pinmyblogs/types/forms"
	"gorm.io/gorm"
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

	h, _ := uuid.NewUUID()
	user.PasswordResetHash = h.String()

	err := database.Db().Save(&user).Error
	if err != nil {
		log.WithFields(map[string]any{
			"email": user.Email,
		}).WithError(result.Error).Errorf("Password reset request %s failed. Failed to update the database.", email)
		c.Redirect(http.StatusServiceUnavailable, "/500")
		return
	}
	resetMailer := mailer.NewResetPasswordMailer(*user)
	go resetMailer.Send()
	c.Redirect(http.StatusSeeOther, "/reset-password/sent?email="+email)
}

func ResetPasswordSentGet(c *gin.Context) {
	log := logger.NewLogger()
	email := c.Query("email")
	if len(email) == 0 {
		log.WithField("email", email).Warn("wrong email in ResetPasswordSentGet")
		c.Redirect(http.StatusFound, "/reset-password")
		c.Abort()
	}

	c.HTML(http.StatusOK, "reset_password_sent.tmpl", gin.H{"Email": email})
}

func ResetPasswordSetGet(c *gin.Context) {
	log := logger.NewLogger()

	hash := c.Param("hash")
	log.Info(hash)
	c.HTML(http.StatusOK, "reset_password_set.tmpl", gin.H{"hash": hash})
}

func ResetPasswordSetPost(c *gin.Context) {
	log := logger.NewLogger()

	field := new(formbinding.FieldErrors)
	form := middlewares.GetForm(c).(*forms.ResetPasswordForm)
	ctx := middlewares.GetContext(c)

	// Helper to show invalid/expired message
	renderInvalidLink := func() {
		c.HTML(http.StatusOK, "reset_password_message.tmpl",
			gin.H{"Message": "This link is expired or invalid."},
		)
	}

	hash := form.Hash
	if hash == "" {
		renderInvalidLink()
		return
	}

	// Fetch user by reset hash
	var user models.User
	if err := database.Db().
		First(&user, "password_reset_hash = ?", hash).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.WithField("passwordResetHash", hash).
				Error("password reset link not found in database")
		} else {
			log.WithError(err).
				Error("database error while fetching user")
		}

		renderInvalidLink()
		return
	}

	password := form.Password
	confirmPassword := form.ConfirmPassword

	// Password validation
	if ctx["Hash_HasError"] == false && !field.IsValid(password) {
		ctx["HasError"] = true
		ctx["Password_HasError"] = true
		ctx["Password_Error"] = field.Error("alpha_dash_dot")
		ctx["ConfirmPassword_HasError"] = false
	}

	// Password match validation
	if (ctx["Password_HasError"] == nil || ctx["Password_HasError"] == false) &&
		password != confirmPassword {
		ctx["HasError"] = true
		ctx["Password_HasError"] = true
		ctx["Password_Error"] = ""
		ctx["ConfirmPassword_HasError"] = true
		ctx["ConfirmPassword_Error"] = field.Error("password_not_match")
	}

	// Stop if validation failed
	if ctx["HasError"] == true {
		c.HTML(http.StatusOK, "reset_password_set.tmpl", ctx)
		return
	}

	// Prevent reusing old password
	if err := utils.CompareBCrypt(user.Password, password); err == nil {
		errMsg := "You cannot use old password."
		ctx["HasError"] = true
		ctx["Error"] = errMsg

		log.WithField("email", user.Email).
			Error(errMsg)

		c.HTML(http.StatusOK, "reset_password_set.tmpl", ctx)
		return
	}

	newPasswdHash := utils.HashifyBCrypt(password)
	// Update password and clear reset hash
	if err := database.Db().
		Model(&models.User{}).
		Where("id = ?", user.ID).
		Updates(map[string]any{
			"password":               newPasswdHash,
			"password_reset_hash":    "",
			"last_password_reset_at": time.Now(),
		}).Error; err != nil {
		log.WithFields(map[string]any{
			"user": user.Email,
			"id":   user.ID,
		}).WithError(err).
			Error("failed to reset the password into database")

		renderInvalidLink()
		return
	}

	log.WithFields(map[string]any{
		"user": user.Email,
		"id":   user.ID,
	}).Info("password has been reset")

	c.HTML(http.StatusOK, "reset_password_message.tmpl",
		gin.H{"Message": "Password has been reset."},
	)
}
