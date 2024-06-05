package public

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func StartGet(c *gin.Context) {
	c.HTML(http.StatusOK, "start.tmpl", nil)
}
