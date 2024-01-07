package home

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Home(c *gin.Context) {
	c.HTML(http.StatusOK, "home.tmpl", nil)
}

func Favourite(c *gin.Context) {
	c.HTML(http.StatusOK, "favourite.tmpl", nil)
}
func Archived(c *gin.Context) {
	c.HTML(http.StatusOK, "archived.tmpl", nil)
}

func Favicon(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
