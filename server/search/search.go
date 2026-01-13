package search

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/puni9869/pinmyblogs/models"
	"github.com/puni9869/pinmyblogs/pkg/database"
	"github.com/puni9869/pinmyblogs/pkg/logger"
	"github.com/puni9869/pinmyblogs/pkg/pagination"
	"github.com/puni9869/pinmyblogs/pkg/utils"
	"github.com/puni9869/pinmyblogs/server/middlewares"
)

func Search(c *gin.Context) {
	log := logger.NewLogger()
	session := sessions.Default(c)
	q := c.DefaultQuery("q", "")
	page, limit := utils.GetPageAndLimit(c)
	if q == "" {
		nextRoute, _ := url.Parse("/home")
		query := nextRoute.Query()
		if page > 1 {
			query.Set("page", strconv.Itoa(page))
		}
		if limit > pagination.DefaultLimit {
			query.Set("limit", strconv.Itoa(limit))
		}
		nextRoute.RawQuery = query.Encode()
		c.Redirect(http.StatusFound, nextRoute.String())
		return
	}
	email, _ := session.Get(middlewares.Userkey).(string)
	p := pagination.Pagination[*models.Url]{Page: page, Limit: limit}
	db := database.Db()
	db.Scopes(pagination.Paginate(&models.Url{}, &p)).
		Where("created_by = ? AND is_active = ? AND is_deleted = ?", email, true, false).
		Where("web_link LIKE ? OR title LIKE ?", "%"+q+"%", "%"+q+"%").
		Order("created_at DESC").
		Find(&p.Items)
	if db.Error != nil {
		log.WithError(db.Error).Error("error in fetching search query")
	}

	var sideNavPref models.Setting
	result := db.Model(models.Setting{}).
		Where("created_by = ? AND action = ? ", email, "sideNav").Find(&sideNavPref)
	if result.Error != nil {
		log.WithError(result.Error).Error("failed to get the preferences on home page")
	}
	log.Infof("getting sideNav prefs %s", sideNavPref.Value)

	// "SideNavCollapse": false  || get from user's settings
	sideNavCollapse := sideNavPref.Value == "hide"
	// p is pagination
	c.HTML(http.StatusOK, "search.html", gin.H{"HasError": false, "Pagination": p, "SearchQuery": q, "Email": email, "SideNavCollapse": sideNavCollapse})
}
