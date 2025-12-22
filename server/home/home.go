package home

import (
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/puni9869/pinmyblogs/models"
	"github.com/puni9869/pinmyblogs/pkg/database"
	"github.com/puni9869/pinmyblogs/pkg/logger"
	"github.com/puni9869/pinmyblogs/pkg/spider"
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
	go spider.ScrapeUrl(&url)

	log.Infof("Requested to add %s in tag: %s ", requestBody.Url, requestBody.Tag)
	c.JSON(http.StatusCreated, gin.H{"Status": "OK", "Message": "Weblink Added."})
}

func Home(c *gin.Context) {
	log := logger.NewLogger()
	session := sessions.Default(c)
	currentlyLoggedIn := session.Get(middlewares.Userkey)
	var urls []models.Url
	db := database.Db()
	result := db.Where("created_by =? and  is_active = ? and is_deleted = ? and is_archived = ?", currentlyLoggedIn.(string), true, false, false).
		Order("id desc").
		Limit(100).
		Find(&urls)
	if result.RowsAffected > 0 {
		log.WithField("resultCount", result.RowsAffected).Info("Fetching the result")
	}

	c.HTML(http.StatusOK, "home.tmpl", gin.H{"HasError": false, "Urls": urls, "Count": result.RowsAffected})
}

func Favourite(c *gin.Context) {
	log := logger.NewLogger()
	session := sessions.Default(c)
	currentlyLoggedIn := session.Get(middlewares.Userkey)
	var urls []models.Url
	db := database.Db()
	result := db.Where("created_by =? and  is_active = ? and is_deleted = ? and is_fav =? ", currentlyLoggedIn.(string), true, false, true).
		Order("id desc").
		Limit(100).
		Find(&urls)
	if result.RowsAffected > 0 {
		log.WithField("resultCount", result.RowsAffected).Info("Fetching the result")
	}

	c.HTML(http.StatusOK, "home.tmpl", gin.H{"HasError": false, "Urls": urls, "Count": result.RowsAffected})
}

func Archived(c *gin.Context) {
	log := logger.NewLogger()
	session := sessions.Default(c)
	currentlyLoggedIn := session.Get(middlewares.Userkey)
	var urls []models.Url
	db := database.Db()
	result := db.Where("created_by =? and  is_active = ? and is_deleted = ? and is_archived =? ", currentlyLoggedIn.(string), true, false, true).
		Order("id desc").
		Limit(100).
		Find(&urls)
	if result.RowsAffected > 0 {
		log.WithField("resultCount", result.RowsAffected).Info("Fetching the result")
	}

	c.HTML(http.StatusOK, "home.tmpl", gin.H{"HasError": false, "Urls": urls, "Count": result.RowsAffected})
}

func Trash(c *gin.Context) {
	log := logger.NewLogger()
	session := sessions.Default(c)
	currentlyLoggedIn := session.Get(middlewares.Userkey)
	var urls []models.Url
	db := database.Db()
	result := db.Where("created_by =? and  is_active = ? and is_deleted = ?", currentlyLoggedIn.(string), true, true).
		Order("updated_at desc").
		Limit(100).
		Find(&urls)
	if result.RowsAffected > 0 {
		log.WithField("resultCount", result.RowsAffected).Info("Fetching the result")
	}

	c.HTML(http.StatusOK, "home.tmpl", gin.H{"HasError": false, "Urls": urls, "Count": result.RowsAffected})
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
