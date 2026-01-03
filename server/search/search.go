package search

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/puni9869/pinmyblogs/models"
	"github.com/puni9869/pinmyblogs/pkg/database"
	"github.com/puni9869/pinmyblogs/pkg/logger"
	"github.com/puni9869/pinmyblogs/pkg/pagination"
	"github.com/puni9869/pinmyblogs/pkg/utils"
	"github.com/puni9869/pinmyblogs/server/middlewares"
	"net/http"
	"net/url"
)

func Search(c *gin.Context) {
	log := logger.NewLogger()
	session := sessions.Default(c)
	q := c.DefaultQuery("q", "")
	page, limit := utils.GetPageAndLimit(c)
	if q == "" {
		nextRoute, _ := url.Parse("/home")
		if page > 1 {
			nextRoute.Query().Set("page", fmt.Sprintf("%d", page))
		}
		if page > pagination.DefaultLimit {
			nextRoute.Query().Set("limit", fmt.Sprintf("%d", limit))
		}

		c.Redirect(http.StatusFound, nextRoute.String())
		return
	}
	email, _ := session.Get(middlewares.Userkey).(string)
	p := pagination.Pagination[*models.Url]{Page: page, Limit: limit}
	db := database.Db()
	db.Scopes(pagination.Paginate(&models.Url{}, &p)).
		Where("created_by = ? AND is_active = ?", email, true).
		Where("web_link LIKE ? OR title LIKE ?", "%"+q+"%", "%"+q+"%").
		Order("created_at DESC").
		Find(&p.Items)
	if db.Error != nil {
		log.WithError(db.Error).Error("error in fetching search query")
	}
	// p is pagination
	c.HTML(http.StatusOK, "search.tmpl", gin.H{"HasError": false, "Pagination": p, "SearchQuery": q})
}
