package setting

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Setting(c *gin.Context) {
	c.HTML(http.StatusOK, "setting.tmpl", nil)
}
