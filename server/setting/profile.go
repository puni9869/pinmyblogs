package setting

import (
	"net/http"
	"slices"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/puni9869/pinmyblogs/models"
	"github.com/puni9869/pinmyblogs/pkg/database"
	"github.com/puni9869/pinmyblogs/pkg/logger"
	"github.com/puni9869/pinmyblogs/server/middlewares"
)

func ProfileAction(c *gin.Context) {
	profileAction := c.Param("action")
	if !slices.Contains([]string{"public", "private"}, profileAction) {
		profileAction = "private"
	}
	log := logger.NewLogger()
	session := sessions.Default(c)
	currentlyLoggedIn := session.Get(middlewares.Userkey)
	email, ok := currentlyLoggedIn.(string)

	if !ok || currentlyLoggedIn == nil || email == "" {
		c.Redirect(http.StatusPermanentRedirect, "/login")
		return
	}

	var user *models.User
	result := database.Db().First(&user, "email = ?", email)
	if result.Error != nil {
		log.WithField("email", email).WithError(result.Error).Error("record not found")
	}
	isPublic := profileAction == "public"
	log.Info(isPublic)
	user.IsProfilePublic = isPublic
	if err := database.Db().Save(user).Error; err != nil {
		log.WithField("profileAction", profileAction).
			WithError(err).
			Error("failed to update profile visibility")
	}

	c.Redirect(http.StatusFound, "/setting")
}
