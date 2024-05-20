package setting

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setting(c *gin.Context) {
	c.HTML(http.StatusOK, "setting.tmpl", nil)
}
