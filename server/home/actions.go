package home

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/puni9869/pinmyblogs/models"
	"github.com/puni9869/pinmyblogs/pkg/database"
	"github.com/puni9869/pinmyblogs/pkg/logger"
	"github.com/puni9869/pinmyblogs/server/middlewares"
	"github.com/puni9869/pinmyblogs/types/forms"
	"net/http"
	"strconv"
)

func BulkActions(c *gin.Context) {
	log := logger.NewLogger()

	session := sessions.Default(c)
	email := session.Get(middlewares.Userkey)

	ctx := middlewares.GetContext(c)
	form := middlewares.GetForm(c).(*forms.BulkAction)
	action := form.Action
	ids := form.IDs
	actions := map[string]string{
		"bulk-delete":    "is_deleted",
		"bulk-archive":   "is_archived",
		"bulk-favourite": "is_fav",
	}
	a, ok := actions[action]
	if ctx["HasError"] == true ||
		ctx["Action_HasError"] == true ||
		ctx["IDs_HasError"] == true || !ok {
		log.WithFields(map[string]any{"ctx": ctx}).Error("bad request for bulk action")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Status": "NOT_OK", "Errors": "Bad values"})
		return
	}

	log.WithFields(map[string]any{
		"email":  email,
		"action": a,
		"ids":    ids,
	}).Info("getting values for bulk action")

	updates := map[string]any{}
	switch action {
	case "bulk-delete":
		updates["is_deleted"] = true
	case "bulk-archive":
		updates["is_archived"] = true
	case "bulk-favourite":
		updates["is_fav"] = true
	}

	db := database.Db()
	if err := db.Model(&models.Url{}).
		Where("id IN ? AND created_by = ?", ids, email).
		Updates(updates).Error; err != nil {

		log.WithFields(map[string]any{"email": email, "action": a, "ids": ids}).
			WithError(err).Error("bulk write to database failed.")

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Status": "NOT_OK", "Errors": "Something went wrong"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Status": "OK", "Message": "Bulk operation successful."})
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
