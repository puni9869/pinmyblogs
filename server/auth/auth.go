package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/puni9869/pinmyblogs/models"
	"github.com/puni9869/pinmyblogs/pkg/database"
)

const userkey = "user"

func LoginPost(c *gin.Context) {
	session := sessions.Default(c)
	username := c.PostForm("email")
	password := c.PostForm("password")

	// Validate form input
	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Parameters can't be empty"})
		return
	}

	// Check for username and password match, usually from a database
	if username != "hello" || password != "itsme" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		return
	}

	// Save the username in the session
	session.Set(userkey, username) // In real world usage you'd set this to the users ID
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, "/home")
}

func LoginGet(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	fmt.Println(user)
	if user == nil {
		c.HTML(http.StatusOK, "login.tmpl", nil)
		c.Abort()
	}
	fmt.Println("Escape Get")
	c.Redirect(http.StatusTemporaryRedirect, "/home")
}

// Logout is the handler called for the user to log out.
func Logout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		c.HTML(http.StatusBadRequest, "index.tmpl", nil)
		return
	}

	session.Delete(userkey)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.HTML(http.StatusOK, "index.tmpl", nil)
	c.Abort()
}

func Signup(c *gin.Context) {
	user := models.User{Name: "Jinzhu", Age: 1}
	db := database.Db()
	db.Create(&user)
	c.HTML(http.StatusOK, "signup.tmpl", nil)
}
