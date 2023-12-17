package home

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Home(c *gin.Context) {
	c.HTML(http.StatusOK, "home.tmpl", gin.H{"title": "Home"})
}

func Favicon(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
