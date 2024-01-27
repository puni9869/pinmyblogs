package auth

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/puni9869/pinmyblogs/models"
	"github.com/puni9869/pinmyblogs/pkg/database"
	"net/http"
)

func Index(c *gin.Context) {
	fmt.Println("INdex")
	c.HTML(http.StatusOK, "index.tmpl", nil)
}

func Login(c *gin.Context) {
	session := sessions.Default(c)
	fmt.Println(session.ID())
	if len(session.ID()) == 0 {
		c.HTML(http.StatusOK, "login.tmpl", nil)
	}
	fmt.Println(session.ID())
	fmt.Println("Here")
	session.Set("isLoggedIn", true)
	session.Set("user", "Hello")
	_ = session.Save()
	c.Redirect(http.StatusTemporaryRedirect, "/home")
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete(session.ID())
	session.Clear()
	fmt.Println(session.ID())
	c.Redirect(http.StatusPermanentRedirect, "/index")
}

func Signup(c *gin.Context) {
	user := models.User{Name: "Jinzhu", Age: 1}
	db := database.Db()
	db.Create(&user)
	c.HTML(http.StatusOK, "signup.tmpl", nil)
}
