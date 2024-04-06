package auth

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

const userkey = "user"

func LoginPost(c *gin.Context) {
	session := sessions.Default(c)
	username := c.PostForm("email")
	password := c.PostForm("password")

	// Validate form input and check authentication
	if username != "hello" || password != "itsme" {
		c.HTML(http.StatusOK, "login.tmpl", nil)
		return
	}
	// Save the username in the session
	session.Set(userkey, username)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	// Redirect to the home route upon successful login
	c.HTML(http.StatusAccepted, "home.tmpl", nil)
}

func LoginGet(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		c.HTML(http.StatusOK, "login.tmpl", nil)
		return
	}
	c.HTML(http.StatusAccepted, "home.tmpl", nil)
}

// Logout is the handler called for the user to log out.
func Logout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		c.Abort()
		return
	}

	session.Delete(userkey)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		c.Abort()
	}
	c.HTML(http.StatusOK, "index.tmpl", nil)
	c.Abort()
}
