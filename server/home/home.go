package home

import (
	"github.com/puni9869/pinmyblogs/pkg/spider"
	"net/http"
	"strconv"

	"github.com/puni9869/pinmyblogs/pkg/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/puni9869/pinmyblogs/models"
	"github.com/puni9869/pinmyblogs/pkg/database"
	"github.com/puni9869/pinmyblogs/pkg/logger"
	"github.com/puni9869/pinmyblogs/pkg/pagination"
	"github.com/puni9869/pinmyblogs/server/middlewares"
	"github.com/puni9869/pinmyblogs/types/forms"
)

func AddWeblink(c *gin.Context) {
	log := logger.NewLogger()
	var err error

	session := sessions.Default(c)
	currentlyLoggedIn := session.Get(middlewares.Userkey)
	requestBody := middlewares.GetForm(c).(*forms.WeblinkRequest)
	ctx := middlewares.GetContext(c)
	if ctx["Tag_HasError"] == true || ctx["Url_HasError"] == true {
		log.WithError(err).Error("Bad request body")
		c.JSON(http.StatusBadRequest, gin.H{"Status": "NOT_OK", "Errors": ctx})
		return
	}

	db := database.Db()
	url := models.Url{WebLink: requestBody.Url,
		IsActive: true, IsDeleted: false,
		CreatedBy: currentlyLoggedIn.(string), Tag: requestBody.Tag,
	}
	db.Save(&url)
	go func(url *models.Url) {
		defer func() {
			if r := recover(); r != nil {
				log.WithFields(map[string]any{
					"panic":   r,
					"webLink": url.WebLink,
				}).Error("Recovered from panic while scraping url")
			}
		}()
		spider.FetchAndUpdateURL(url)
	}(&url)
	log.Infof("Requested to add %s in tag: %s ", requestBody.Url, requestBody.Tag)
	c.JSON(http.StatusCreated, gin.H{"Status": "OK", "Message": "Weblink Added."})
}

func Home(c *gin.Context) {
	log := logger.NewLogger()
	session := sessions.Default(c)
	email, _ := session.Get(middlewares.Userkey).(string)
	page, limit := utils.GetPageAndLimit(c)
	p := pagination.Pagination[*models.Url]{Page: page, Limit: limit}
	db := database.Db()
	db.Scopes(pagination.Paginate(&models.Url{}, &p)).
		Where(
			"created_by = ? AND is_active = ? AND is_deleted = ? AND is_archived = ?",
			email, true, false, false,
		).
		Order("id DESC").
		Find(&p.Items)

	var sideNavPref models.Setting
	result := database.Db().Model(models.Setting{}).
		Where("created_by = ? AND action = ? ", email, "sideNav").Find(&sideNavPref)
	if result.Error != nil {
		log.WithError(result.Error).Error("failed to get the preferences on home page")
	}
	log.Infof("getting sideNav prefs %s", sideNavPref.Value)

	// "SideNavCollapse": false  || get from user's settings
	sideNavCollapse := sideNavPref.Value == "hide"
	// p is pagination
	c.HTML(http.StatusOK, "home.tmpl", gin.H{
		"HasError":        false,
		"Pagination":      p,
		"Email":           email,
		"SideNavCollapse": sideNavCollapse,
	})
}

func Favourite(c *gin.Context) {
	log := logger.NewLogger()
	session := sessions.Default(c)
	email, _ := session.Get(middlewares.Userkey).(string)
	page, limit := utils.GetPageAndLimit(c)
	p := pagination.Pagination[*models.Url]{Page: page, Limit: limit}
	db := database.Db()
	db.Scopes(pagination.Paginate(&models.Url{}, &p)).
		Where("created_by =? and  is_active = ? and is_deleted = ? and is_fav =? ", email, true, false, true).
		Order("id desc").
		Find(&p.Items)
	var sideNavPref models.Setting
	result := database.Db().Model(models.Setting{}).
		Where("created_by = ? AND action = ? ", email, "sideNav").Find(&sideNavPref)
	if result.Error != nil {
		log.WithError(result.Error).Error("failed to get the preferences on home page")
	}
	log.Infof("getting sideNav prefs %s", sideNavPref.Value)

	// "SideNavCollapse": false  || get from user's settings
	sideNavCollapse := sideNavPref.Value == "hide"
	// p is pagination
	c.HTML(http.StatusOK, "home.tmpl", gin.H{
		"HasError":        false,
		"Pagination":      p,
		"Email":           email,
		"SideNavCollapse": sideNavCollapse,
	})
}

func Archived(c *gin.Context) {
	log := logger.NewLogger()
	session := sessions.Default(c)
	email, _ := session.Get(middlewares.Userkey).(string)
	page, limit := utils.GetPageAndLimit(c)
	p := pagination.Pagination[*models.Url]{Page: page, Limit: limit}
	db := database.Db()
	db.Scopes(pagination.Paginate(&models.Url{}, &p)).
		Where("created_by =? and  is_active = ? and is_deleted = ? and is_archived =? ", email, true, false, true).
		Order("id desc").
		Find(&p.Items)
	var sideNavPref models.Setting
	result := database.Db().Model(models.Setting{}).
		Where("created_by = ? AND action = ? ", email, "sideNav").Find(&sideNavPref)
	if result.Error != nil {
		log.WithError(result.Error).Error("failed to get the preferences on home page")
	}
	log.Infof("getting sideNav prefs %s", sideNavPref.Value)

	// "SideNavCollapse": false  || get from user's settings
	sideNavCollapse := sideNavPref.Value == "hide"
	// p is pagination
	c.HTML(http.StatusOK, "home.tmpl", gin.H{
		"HasError":        false,
		"Pagination":      p,
		"Email":           email,
		"SideNavCollapse": sideNavCollapse,
	})
	c.HTML(http.StatusOK, "home.tmpl", gin.H{"HasError": false, "Pagination": p, "Email": email})
}

func Trash(c *gin.Context) {
	log := logger.NewLogger()
	session := sessions.Default(c)
	email, _ := session.Get(middlewares.Userkey).(string)
	page, limit := utils.GetPageAndLimit(c)
	p := pagination.Pagination[*models.Url]{Page: page, Limit: limit}
	db := database.Db()
	db.Scopes(pagination.Paginate(&models.Url{}, &p)).
		Where("created_by =? and  is_active = ? and is_deleted = ?", email, true, true).
		Order("id desc").
		Find(&p.Items)
	var sideNavPref models.Setting
	result := database.Db().Model(models.Setting{}).
		Where("created_by = ? AND action = ? ", email, "sideNav").Find(&sideNavPref)
	if result.Error != nil {
		log.WithError(result.Error).Error("failed to get the preferences on home page")
	}
	log.Infof("getting sideNav prefs %s", sideNavPref.Value)

	// "SideNavCollapse": false  || get from user's settings
	sideNavCollapse := sideNavPref.Value == "hide"
	// p is pagination
	c.HTML(http.StatusOK, "home.tmpl", gin.H{
		"HasError":        false,
		"Pagination":      p,
		"Email":           email,
		"SideNavCollapse": sideNavCollapse,
	})
}

func Actions(c *gin.Context) {
	log := logger.NewLogger()
	var err error
	session := sessions.Default(c)
	currentlyLoggedIn := session.Get(middlewares.Userkey)

	var requestBody map[string]string
	err = c.ShouldBindJSON(&requestBody)
	if err != nil {
		log.WithError(err).Error("Bad request body")
		c.JSON(http.StatusBadRequest, gin.H{"Status": "NOT_OK", "Errors": "Bad values"})
		return
	}

	u64, err := strconv.ParseUint(requestBody["id"], 10, 64)
	if err != nil {
		log.WithError(err).Error("Bad request body", u64)
		c.JSON(http.StatusBadRequest, gin.H{"Status": "NOT_OK", "Errors": "Bad values"})
		return
	}

	var updates = make(map[string]any)

	if val, ok := requestBody["isFav"]; ok {
		value, err := strconv.ParseBool(val)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Status": "NOT_OK", "Errors": "Bad values"})
			return
		}
		updates["IsFav"] = !value
	}

	if val, ok := requestBody["isArchived"]; ok {
		value, err := strconv.ParseBool(val)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Status": "NOT_OK", "Errors": "Bad values"})
			return
		}
		updates["IsArchived"] = !value
	}

	if val, ok := requestBody["isDeleted"]; ok {
		value, err := strconv.ParseBool(val)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Status": "NOT_OK", "Errors": "Bad values"})
		}
		updates["IsDeleted"] = !value
	}

	db := database.Db()
	db.Model(&models.Url{}).Where("id = ? and created_by = ? ", requestBody["id"], currentlyLoggedIn.(string)).Updates(updates)

	c.JSON(http.StatusOK, gin.H{"Status": "OK", "Message": "Weblink updated."})
}

func Share(c *gin.Context) {
	log := logger.NewLogger()
	//var err error
	session := sessions.Default(c)
	currentlyLoggedIn := session.Get(middlewares.Userkey)
	id := c.Param("id")
	var url *models.Url
	db := database.Db()
	result := db.Where("id =? and created_by =? and  is_active = ?", id, currentlyLoggedIn.(string), true).First(&url)
	if result.RowsAffected != 1 {
		c.JSON(http.StatusNotFound, map[string]string{"Status": "NOT_OK", "Message": "Not found."})
		return
	}
	log.Info(url)
	c.HTML(http.StatusOK, "share.tmpl", nil)
}
