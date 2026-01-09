package setting

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/puni9869/pinmyblogs/models"
	"github.com/puni9869/pinmyblogs/pkg/database"
	"github.com/puni9869/pinmyblogs/pkg/logger"
	"github.com/puni9869/pinmyblogs/server/middlewares"
	"github.com/puni9869/pinmyblogs/types/forms"
	"gorm.io/gorm/clause"
	"net/http"
	"slices"
)

var preferences = map[string][]string{
	"sideNav": {"show", "hide"},
}

func Prefs(c *gin.Context) {
	log := logger.NewLogger()
	form := middlewares.GetForm(c).(*forms.Prefs)
	ctx := middlewares.GetContext(c)
	action := form.Action
	value := form.Value

	session := sessions.Default(c)
	email, _ := session.Get(middlewares.Userkey).(string)

	logFields := map[string]any{
		"email":  email,
		"action": action,
		"value":  value,
	}
	log = log.WithFields(logFields).Logger
	if values, ok := preferences[action]; !ok || !slices.Contains(values, value) {
		log.Error("invalid action and value")
		c.AbortWithStatusJSON(http.StatusBadRequest, ctx)
		return
	}
	db := database.Db()
	s := models.Setting{}
	s.Value = value
	s.Action = action
	s.CreatedBy = email

	err := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "created_by"},
			{Name: "action"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"value",
			"categories",
			"is_show_count",
			"updated_at",
		}),
	}).Create(&s).Error
	if err != nil {
		log.WithError(err).Error("failed to update the user's setting prefs. database error.")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	log.Info("prefs saved to database")
	c.AbortWithStatus(http.StatusOK)
}
