package setting

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/puni9869/pinmyblogs/models"
	"github.com/puni9869/pinmyblogs/pkg/config"
	"github.com/puni9869/pinmyblogs/pkg/database"
	"github.com/puni9869/pinmyblogs/pkg/logger"
	"github.com/puni9869/pinmyblogs/server/middlewares"

	"github.com/gin-gonic/gin"
)

func Setting(c *gin.Context) {
	log := logger.NewLogger()
	session := sessions.Default(c)
	currentlyLoggedIn := session.Get(middlewares.Userkey)
	tmplCtx := make(map[string]any)
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
		tmplCtx["IsProfilePublic"] = user.IsProfilePublic
	}
	var sideNavPref models.Setting
	result := database.Db().Model(models.Setting{}).
		Where("created_by = ? AND action = ? ", email, "sideNav").Find(&sideNavPref)
	if result.Error != nil {
		log.WithError(result.Error).Error("failed to get the preferences on home page")
	}
	log.Infof("getting sideNav prefs %s", sideNavPref.Value)
	// "SideNavCollapse": false  || get from user's settings
	tmplCtx["SideNavCollapse"] = sideNavPref.Value == "hide"
	c.HTML(http.StatusOK, "setting.html", tmplCtx)
}

func DownloadMyData(c *gin.Context) {
	log := logger.NewLogger()

	format := strings.ToLower(c.Param("format"))
	if !slices.Contains([]string{"csv", "json", "html"}, format) {
		format = "json"
	}
	log.Infof("generating the data in %s format", format)
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
