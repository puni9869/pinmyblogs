package auth

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func LoginPost(c *gin.Context) {
	session := sessions.Default(c)
	username := c.PostForm("email")
	password := c.PostForm("password")

	// Validate form input and check authentication
	if username != "hello" || password != "itsme" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	// Save the username in the session
	session.Set(userkey, username)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	fmt.Println(c.Request.URL.Path)
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
	fmt.Println("Escape Get")
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
		return
	}
	c.HTML(http.StatusOK, "index.tmpl", nil)
	c.Abort()
	return
}
