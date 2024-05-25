package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/puni9869/pinmyblogs/pkg/logger"
	"github.com/puni9869/pinmyblogs/server/middlewares"
	"github.com/puni9869/pinmyblogs/types/forms"
)

func ResetPasswordGet(c *gin.Context) {
	c.HTML(http.StatusOK, "reset.tmpl", nil)
}

func ResetPasswordPost(c *gin.Context) {
	log := logger.NewLogger()

	form := middlewares.GetForm(c).(*forms.ResetForm)
	ctx := middlewares.GetContext(c)

	email := form.Email

	if ctx["Email_HasError"] == true {
		fmt.Print("Here")
		log.WithField("email", email).Error("email id not found.")
		c.HTML(http.StatusBadRequest, "reset.tmpl", gin.H{"Email": email, "HasError": true, "Error": "Email id not found."})
		return
	}

	c.HTML(http.StatusAccepted, "reset_password_sent.tmpl", gin.H{"Email": email})
}
