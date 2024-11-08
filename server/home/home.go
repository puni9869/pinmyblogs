package home

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/puni9869/pinmyblogs/pkg/logger"
)

func AddWeblink(c *gin.Context) {
	log := logger.NewLogger()

	var requestBody struct {
		Url string `json:"url"`
		Tag string `json:"tag"`
	}
	// Handle the error for url validation
	_ = c.ShouldBindJSON(&requestBody)
	log.Info("Requested to add %s in tag: %s ", requestBody.Url, requestBody.Tag)
	c.JSON(http.StatusCreated, gin.H{"Status": "OK", "Message": "Weblink Added."})
}

func Home(c *gin.Context) {
	c.HTML(http.StatusOK, "home.tmpl", nil)
}

func Favourite(c *gin.Context) {
	c.HTML(http.StatusOK, "favourite.tmpl", nil)
}

func Archived(c *gin.Context) {
	c.HTML(http.StatusOK, "archived.tmpl", nil)
}

func Trash(c *gin.Context) {
	c.HTML(http.StatusOK, "trash.tmpl", nil)
}

func Favicon(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"status": "OK"})
}
func OK(c *gin.Context) {
	log := logger.NewLogger()
	p := c.Request.URL
	log.Infof("%#v", p)
	c.String(http.StatusOK, "OK")
}
