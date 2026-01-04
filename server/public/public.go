package public

import (
	"crypto/sha256"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/puni9869/pinmyblogs"
)

var (
	favicon     []byte
	faviconETag string
)

func LandingPageGet(c *gin.Context) {
	c.HTML(http.StatusOK, "landing_page.tmpl", nil)
}

func LandingPage2Get(c *gin.Context) {
	c.HTML(http.StatusOK, "landing_page_2.tmpl", nil)
}

func PrivacyPolicyGet(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=31536000, immutable")
	c.Header("ETag", faviconETag)
	c.HTML(http.StatusOK, "privacy_policy.tmpl", nil)
}

func SupportGet(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=31536000, immutable")
	c.Header("ETag", faviconETag)
	c.HTML(http.StatusOK, "support.tmpl", nil)
}

func FavIcon(c *gin.Context) {
	if len(favicon) > 0 {
		c.Data(http.StatusOK, "image/x-icon", favicon)
	}
	var err error
	favicon, err = pinmyblogs.Files.ReadFile("frontend/icons/favicon.ico")
	if err != nil {
		panic("favicon.ico not found")
	}

	hash := sha256.Sum256(favicon)
	faviconETag = fmt.Sprintf(`"%x"`, hash[:8])

	if c.GetHeader("If-None-Match") == faviconETag {
		c.Status(http.StatusNotModified)
		return
	}

	c.Header("Cache-Control", "public, max-age=31536000, immutable")
	c.Header("ETag", faviconETag)
	c.Data(http.StatusOK, "image/x-icon", favicon)
}

func Route404(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404.tmpl", nil)
}

func Route500(c *gin.Context) {
	c.HTML(http.StatusInternalServerError, "500.tmpl", nil)
}

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"status": "OK"})
}

func ServiceWorker(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=31536000, immutable")
	c.Header("ETag", faviconETag)
	c.File("frontend/service-worker.js")
}

func OfflinePage(c *gin.Context) {
	c.HTML(http.StatusOK, "offline_page.tmpl", nil)
}
