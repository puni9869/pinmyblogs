package auth

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/puni9869/pinmyblogs/models"
	"github.com/puni9869/pinmyblogs/pkg/database"
	"github.com/puni9869/pinmyblogs/pkg/logger"
	"net/http"
)

const userkey = "user"

func LoginPost(c *gin.Context) {
	log := logger.NewLogger()

	email := c.PostForm("email")
	password := c.PostForm("password")

	if email == "" || password == "" {
		log.WithField("email", email).Error("Empty password or email.")
		c.HTML(http.StatusBadRequest, "login.tmpl", gin.H{"HasError": true, "Error": "Invalid email or password."})
		c.Abort()
	}

	var user *models.User
	result := database.Db().First(&user, "email = ?", email)
	if result.Error != nil {
		log.WithField("email", email).WithError(result.Error).Error("Invalid email or password. Database error")
		c.HTML(http.StatusUnauthorized, "login.tmpl", gin.H{"HasError": true, "Error": "Invalid email or password"})
		c.Abort()
	}

	// Save the username in the session
	session := sessions.Default(c)
	currentlyLoggedIn := session.Get(userkey)
	log.WithField("email", email).Info("login found ", currentlyLoggedIn)
	if currentlyLoggedIn == nil || currentlyLoggedIn != email {
		session.Set(userkey, email)
		log.WithField("email", email).Info("setting user", currentlyLoggedIn)
		if err := session.Save(); err != nil {
			log.WithField("email", email).WithError(result.Error).Error("Unable to save session.")
			c.HTML(http.StatusInternalServerError, "login.tmpl", gin.H{"HasError": true, "Error": "Something went wrong. We are working on it."})
			c.Abort()
		}
		log.WithField("email", email).Info("set user", currentlyLoggedIn)
	}
	fmt.Println(user, "user currently logged in")
	// Redirect to the home route upon successful login
	c.HTML(http.StatusOK, "home.tmpl", nil)
}

func LoginGet(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tmpl", nil)
	// session := sessions.Default(c)
	// user := session.Get(userkey)
	//
	//	if user == nil {
	//		c.HTML(http.StatusOK, "login.tmpl", nil)
	//		return
	//	}
	//
	// c.HTML(http.StatusAccepted, "home.tmpl", nil)
}

// Logout is the handler called for the user to log out.
func Logout(c *gin.Context) {
	log := logger.NewLogger()

	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		log.WithField("user", user).Info("Redirecting to login page. Session not found")
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		c.Abort()
		return
	}
	sessionId := session.ID()
	log.Info("session id ", sessionId)

	session.Clear()
	session.Delete(userkey)
	if len(sessionId) != 0 {
		log.WithField("session", user).Info("session id found")
		res := database.Db().Table("sessions").Where("id = ?", sessionId)
		log.Info("rows affected ", res.RowsAffected)
		if res.Error != nil {
			log.WithField("session", user).WithError(res.Error).Error("failed to delete the session")
		}
	}
	c.Redirect(http.StatusTemporaryRedirect, "/login")
}
