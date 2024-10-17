package setting

import (
	"github.com/puni9869/pinmyblogs/models"
	"github.com/puni9869/pinmyblogs/pkg/database"
	"github.com/puni9869/pinmyblogs/pkg/logger"
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/puni9869/pinmyblogs/server/middlewares"

	"github.com/gin-gonic/gin"
)

func Setting(c *gin.Context) {
	log := logger.NewLogger()
	session := sessions.Default(c)
	currentlyLoggedIn := session.Get(middlewares.Userkey)
	tmplCtx := make(map[string]string)
	email := currentlyLoggedIn.(string)
	if currentlyLoggedIn != nil && len(email) > 0 {
		tmplCtx["Email"] = currentlyLoggedIn.(string)
		var user *models.User
		result := database.Db().First(&user, "email = ?", email)
		if result.Error != nil {
			log.WithField("email", email).WithError(result.Error).Error("record not found")
		}
		if user.DisplayName == "" {
			tmplCtx["DisplayName"] = strings.Split(email, "@")[0]
		} else {
			tmplCtx["DisplayName"] = user.DisplayName
		}

	}
	c.HTML(http.StatusOK, "setting.tmpl", tmplCtx)
}

func DeleteMyAccount(c *gin.Context) {
	c.HTML(http.StatusOK, "settings_page.tmpl", nil)
}

func DisableMyAccount(c *gin.Context) {
	log := logger.NewLogger()

	session := sessions.Default(c)
	email := session.Get(middlewares.Userkey)
	log.Infof("Got disablemyaccount request for user-%s \n", email)

	var user *models.User
	result := database.Db().First(&user, "email = ?", email)
	if result.Error != nil {
		log.WithField("email", nil).WithError(result.Error).Error("Email not found. Database error")
		c.JSON(http.StatusNotFound, map[string]string{"Status": "NOT_OK", "Message": "Account not found."})
		c.Abort()
		return
	}
	if user == nil {
		log.WithField("email", nil).WithError(result.Error).Error("User not found and getting nil from database.")
		c.JSON(http.StatusNotFound, map[string]string{"Status": "NOT_OK", "Message": "User not found."})
		c.Abort()
		return
	}
	user.IsActive = false
	database.Db().Save(user)
	c.JSON(http.StatusOK, map[string]string{"Status": "OK"})
}
