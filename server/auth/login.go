package auth

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/puni9869/pinmyblogs/models"
	"github.com/puni9869/pinmyblogs/pkg/config"
	"github.com/puni9869/pinmyblogs/pkg/database"
	"github.com/puni9869/pinmyblogs/pkg/logger"
	"github.com/puni9869/pinmyblogs/pkg/utils"
)

const Userkey = "user"

func LoginPost(c *gin.Context) {
	log := logger.NewLogger()
	email := c.PostForm("email")
	password := c.PostForm("password")

	if !config.C.Authentication.EnableLogin {
		log.WithFields(map[string]any{
			"enableLoginFlag": config.C.Authentication.EnableLogin,
			"email":           email,
		}).Error("Login service is globally disabled.")
		c.HTML(http.StatusForbidden, "login.tmpl",
			gin.H{"HasError": true, "Error": "Account's login is currently disabled. We are working on it."},
		)
		c.Abort()
		return
	}

	if email == "" || password == "" || !utils.ValidateEmail(email) {
		log.WithField("email", email).Error("Empty password or email. Or email address is not valid.")
		c.HTML(http.StatusBadRequest, "login.tmpl", gin.H{"HasError": true, "Error": "Invalid email or password."})
		c.Abort()
	}

	var user *models.User
	result := database.Db().First(&user, "email = ?", email)
	if result.Error != nil {
		log.WithField("email", email).WithError(result.Error).Error("Invalid email or password. Database error")
		c.HTML(http.StatusUnauthorized, "login.tmpl", gin.H{"HasError": true, "Error": "Account not found."})
		c.Abort()
		return
	}

	if config.C.Authentication.OpenDisabledAccountByEmailLink {
		if !user.IsActive {
			log.WithFields(map[string]any{
				"email":           user.Email,
				"isActive":        user.IsActive,
				"isEmailVerified": user.EmailVerifyHash,
			}).WithError(result.Error).Error("Account is disabled.")
			c.HTML(http.StatusUnauthorized, "login.tmpl", gin.H{"HasError": true, "Error": "Account is disabled. Please check your email."})
			c.Abort()
			return
		}
	}

	if !user.IsEmailVerified {
		log.WithFields(map[string]any{
			"email":           user.Email,
			"isActive":        user.IsActive,
			"isEmailVerified": user.EmailVerifyHash,
		}).WithError(result.Error).Error("Account is not verified.")
		c.HTML(http.StatusUnauthorized, "login.tmpl", gin.H{"HasError": true, "Error": "Account is not verified. Please check your email."})
		c.Abort()
		return
	}
	passwordHash := utils.HashifySHA256(password)
	if strings.Compare(passwordHash, user.Password) != 0 {
		log.WithField("email", email).WithError(result.Error).Error("Invalid password.")
		c.HTML(http.StatusUnauthorized, "login.tmpl", gin.H{"HasError": true, "Error": "Invalid email or password"})
		c.Abort()
		return
	}

	// Save the username in the session
	session := sessions.Default(c)
	currentlyLoggedIn := session.Get(Userkey)
	log.WithField("email", email).Info("login found ", currentlyLoggedIn)

	if currentlyLoggedIn == nil || currentlyLoggedIn != email {
		session.Set(Userkey, email)
		log.WithField("email", email).Info("setting user", currentlyLoggedIn)

		if err := session.Save(); err != nil {
			log.WithField("email", email).WithError(result.Error).Error("Unable to save session.")
			c.HTML(http.StatusInternalServerError, "login.tmpl", gin.H{"HasError": true, "Error": "Something went wrong. We are working on it."})
			c.Abort()
		}

		log.WithField("email", email).Info("set user", currentlyLoggedIn)
	}

	log.WithField("user", user).Info("user currently logged in")
	// Redirect to the home route upon successful login
	c.Redirect(http.StatusPermanentRedirect, "/")
	c.Abort()
}

func LoginGet(c *gin.Context) {
	log := logger.NewLogger()

	session := sessions.Default(c)
	currentlyLoggedIn := session.Get(Userkey)

	if currentlyLoggedIn == nil || len(currentlyLoggedIn.(string)) == 0 {
		c.HTML(http.StatusOK, "login.tmpl", nil)
		c.Abort()
		return
	}
	var user *models.User
	result := database.Db().First(&user, "email = ?", currentlyLoggedIn)
	if result.Error != nil {
		log.WithField("email", currentlyLoggedIn).WithError(result.Error).Error("User not found in database. Database error")
		c.HTML(http.StatusUnauthorized, "login.tmpl", gin.H{"HasError": true, "Error": "Invalid email or password"})
		c.Abort()
		return
	}
	log.WithField("email", currentlyLoggedIn).Info("loggedIn user")
	c.HTML(http.StatusOK, "home.tmpl", nil)
	c.Abort()
}

// Logout is the handler called for the user to log out.
func Logout(c *gin.Context) {
	log := logger.NewLogger()

	session := sessions.Default(c)
	user := session.Get(Userkey)
	if user == nil {
		log.WithField("user", user).Info("Redirecting to login page. Session not found")
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		c.Abort()
	}

	sId := session.ID()

	if len(sId) != 0 {
		log.WithField("session", user).Info("session id found")

		session.Delete(sId)
		session.Set(Userkey, nil)

		var s *models.Session
		res := database.Db().Delete(&s, "id = ?", sId)
		if res.Error != nil {
			log.WithField("session", user).WithError(res.Error).Error("failed to delete the session from database")
		}

		log.Info("rows affected ", res.RowsAffected)

		if err := session.Save(); err != nil {
			log.WithError(err).Error("Unable to delete the session.")
			c.HTML(http.StatusInternalServerError, "login.tmpl", gin.H{"HasError": true, "Error": "Something went wrong. We are working on it."})
			c.Abort()
		}
	}
	c.Redirect(http.StatusTemporaryRedirect, "/login")
}
