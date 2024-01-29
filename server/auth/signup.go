package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/puni9869/pinmyblogs/models"
	"github.com/puni9869/pinmyblogs/pkg/database"
)

func Signup(c *gin.Context) {
	user := models.User{FirstName: "Matt", LastName: "R"}
	db := database.Db()
	db.Create(&user)
	c.HTML(http.StatusOK, "signup.tmpl", nil)
}
