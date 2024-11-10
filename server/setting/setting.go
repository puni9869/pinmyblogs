package setting

import (
	"fmt"
	"github.com/puni9869/pinmyblogs/models"
	"github.com/puni9869/pinmyblogs/pkg/config"
	"github.com/puni9869/pinmyblogs/pkg/database"
	"github.com/puni9869/pinmyblogs/pkg/logger"
	"net/http"
	"slices"
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
		tmplCtx["ShareDataOverMail"] = fmt.Sprintf("%t", config.C.ShareDataOverMail)
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

func DownloadMyData(c *gin.Context) {
	log := logger.NewLogger()

	format := strings.ToLower(c.Param("format"))
	if !slices.Contains([]string{"csv", "json", "html"}, format) {
		format = "json"
	}

	session := sessions.Default(c)
	currentlyLoggedIn := session.Get(middlewares.Userkey)

	var urls []models.Url
	db := database.Db()
	result := db.Where("created_by =? and  is_active = ?", currentlyLoggedIn.(string), true).Find(&urls)
	if result.Error != nil {
		log.WithError(result.Error).Error("Error in fetching the data")
		c.JSON(http.StatusInternalServerError, map[string]string{"Status": "NOT_OK", "Message": "Something went wrong."})
		return
	}
	if result.RowsAffected == 0 {
		log.WithField("resultCount", result.RowsAffected).Info("Fetching the result")
		c.JSON(http.StatusNotFound, map[string]string{"Status": "NOT_OK", "Message": "No Data."})
	}
	c.JSON(http.StatusAccepted, urls)
}
