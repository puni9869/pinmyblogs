package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/puni9869/pinmyblogs/internal/signup"
	"github.com/puni9869/pinmyblogs/pkg/logger"
	"net/http"
)

func SignupGet(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.tmpl", nil)
}

func SignupPost(signupService signup.SignupService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error
		log := logger.NewLogger()
		err = c.Request.ParseForm()
		if err != nil {
			log.WithError(err).Error("Error parsing form data.")
			c.JSON(http.StatusBadRequest, "Error parsing form data.")
			return
		}

		signupService.ValidateForm(c.Request)

		// This will infer what binder to use depending on the content-type header.
		//If err = c.ShouldBind(&signup); err != nil {
		//	log.Error(err)
		//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		//	return
		//}

		//user := models.User{FirstName: "Matt", LastName: "R"}
		//db := database.Db()
		//db.Create(&user)
		c.HTML(http.StatusOK, "signup.tmpl", nil)
	}
}
