package setting

import (
	"errors"
	"net/http"
	"slices"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/puni9869/pinmyblogs/models"
	"github.com/puni9869/pinmyblogs/pkg/config"
	"github.com/puni9869/pinmyblogs/pkg/database"
	"github.com/puni9869/pinmyblogs/pkg/logger"
	"github.com/puni9869/pinmyblogs/pkg/mailer"
	"github.com/puni9869/pinmyblogs/server/middlewares"
	"gorm.io/gorm"
)

func DeleteMyAccount(c *gin.Context) {
	c.HTML(http.StatusOK, "settings_page.tmpl", nil)
}

func DisableMyAccount(c *gin.Context) {
	log := logger.NewLogger()

	session := sessions.Default(c)
	email := session.Get(middlewares.Userkey)
	log.Infof("Got disablemyaccount request for user-%s \n", email)

	var user *models.User
	result := database.Db().First(&user, "email = ?", email)
	if result.Error != nil {
		log.WithField("email", nil).WithError(result.Error).Error("Email not found. Database error")
		c.JSON(http.StatusNotFound, map[string]string{"Status": "NOT_OK", "Message": "Account not found."})
		c.Abort()
		return
	}
	if user == nil {
		log.WithField("email", nil).WithError(result.Error).Error("User not found and getting nil from database.")
		c.JSON(http.StatusNotFound, map[string]string{"Status": "NOT_OK", "Message": "User not found."})
		c.Abort()
		return
	}
	if config.C.Authentication.OpenDisabledAccountByEmailLink {
		h, _ := uuid.NewUUID()
		user.AccountEnableHash = h.String()
		user.IsActive = false
		database.Db().Save(user)

		log.WithField("email", email).Info("user is deactivated")

		action := "disable"
		m := mailer.NewAccountService(*user, action)
		go m.Send()
	}
	c.JSON(http.StatusOK, map[string]string{"Status": "OK"})
}

func EnableMyAccount(c *gin.Context) {
	log := logger.NewLogger()

	hash := c.Param("hash")
	if hash == "" {
		c.HTML(http.StatusOK, "account_message.tmpl",
			gin.H{"Message": "Invalid or expired link."},
		)
		return
	}

	// Validate UUID format early (important)
	if _, err := uuid.Parse(hash); err != nil {
		log.WithField("hash", hash).
			Warn("invalid UUID in enable account link")

		c.HTML(http.StatusOK, "account_message.tmpl",
			gin.H{"Message": "Invalid or expired link."},
		)
		return
	}

	var user models.User
	if err := database.Db().
		Where("account_enable_hash = ?", hash).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.WithField("hash", hash).
				Warn("enable account link not found")
		} else {
			log.WithError(err).
				Error("database error while enabling account")
		}

		c.HTML(http.StatusOK, "account_message.tmpl",
			gin.H{"Message": "This link is invalid or has already been used."},
		)
		return
	}

	// Already enabled safeguard
	if user.IsActive {
		c.HTML(http.StatusOK, "account_message.tmpl",
			gin.H{"Message": "Your account is already active."},
		)
		return
	}

	// Enable account & invalidate hash
	now := time.Now()
	if err := database.Db().
		Model(&models.User{}).
		Where("id = ?", user.ID).
		Updates(map[string]any{
			"is_active":              true,
			"account_enable_hash":    "",
			"last_account_enable_at": &now,
		}).Error; err != nil {
		log.WithFields(map[string]any{
			"user": user.Email,
			"id":   user.ID,
		}).WithError(err).
			Error("failed to enable account")

		c.HTML(http.StatusOK, "account_message.tmpl",
			gin.H{"Message": "Something went wrong. Please try again."},
		)
		return
	}

	log.WithFields(map[string]any{
		"user": user.Email,
		"id":   user.ID,
	}).Info("account successfully enabled")

	c.HTML(http.StatusOK, "account_message.tmpl",
		gin.H{"Message": "Your account has been successfully enabled. You can now log in."},
	)
}

func ProfileAction(c *gin.Context) {
	profileAction := c.Param("action")
	if !slices.Contains([]string{"public", "private"}, profileAction) {
		profileAction = "private"
	}
	log := logger.NewLogger()
	session := sessions.Default(c)
	currentlyLoggedIn := session.Get(middlewares.Userkey)
	email, ok := currentlyLoggedIn.(string)

	if !ok || currentlyLoggedIn == nil || email == "" {
		c.Redirect(http.StatusPermanentRedirect, "/login")
		return
	}

	var user *models.User
	result := database.Db().First(&user, "email = ?", email)
	if result.Error != nil {
		log.WithField("email", email).WithError(result.Error).Error("record not found")
	}
	isPublic := profileAction == "public"
	log.Info(isPublic)
	user.IsProfilePublic = isPublic
	if err := database.Db().Save(user).Error; err != nil {
		log.WithField("profileAction", profileAction).
			WithError(err).
			Error("failed to update profile visibility")
	}

	c.Redirect(http.StatusFound, "/setting")
}
