package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/puni9869/pinmyblogs/pkg/logger"
	"github.com/puni9869/pinmyblogs/types/forms"
	"net/http"
)

type SignUp interface {
	Register()
	Validate()
	CheckPassword()
	CheckEmail()
	IsActive()
	Verify()
	IsVerified()
}

func SignupGet(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.tmpl", nil)
}

func SignupPost(c *gin.Context) {
	log := logger.NewLogger()
	var err error
	err = c.Request.ParseForm()
	if err != nil {
		log.WithError(err).Error("Error parsing form data.")
		c.JSON(http.StatusBadRequest, "Error parsing form data.")
		return
	}

	var signup forms.SignUpForm
	// This will infer what binder to use depending on the content-type header.
	if err = c.ShouldBind(&signup); err != nil {
		log.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//user := models.User{FirstName: "Matt", LastName: "R"}
	//db := database.Db()
	//db.Create(&user)
	c.HTML(http.StatusOK, "signup.tmpl", nil)
}
