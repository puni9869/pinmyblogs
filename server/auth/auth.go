package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/puni9869/pinmyblogs/models"
	"github.com/puni9869/pinmyblogs/pkg/database"
	"net/http"
)

func Login(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tmpl", nil)
}

func Logout(c *gin.Context) {
	c.Redirect(http.StatusFound, "/login")
}

func Signup(c *gin.Context) {
	user := models.User{Name: "Jinzhu", Age: 1}
	db := database.Db()
	db.Create(&user)
	c.HTML(http.StatusOK, "signup.tmpl", nil)
}
