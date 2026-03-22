// Package public provides handlers for unauthenticated public pages.
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

// LandingPageGet renders the public landing page.
func LandingPageGet(c *gin.Context) {
	c.HTML(http.StatusOK, "landing_page.html", nil)
}

// PrivacyPolicyGet renders the privacy policy page.
func PrivacyPolicyGet(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=31536000, immutable")
	c.Header("ETag", faviconETag)
	c.HTML(http.StatusOK, "privacy_policy.html", nil)
}

// SupportGet renders the support page.
func SupportGet(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=31536000, immutable")
	c.Header("ETag", faviconETag)
	c.HTML(http.StatusOK, "support.html", nil)
}

// FavIcon serves the embedded favicon with caching headers.
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

// Route404 renders the 404 not found page.
func Route404(c *gin.Context) {
	c.HTML(http.StatusNotFound, "404.html", nil)
}

// Route500 renders the 500 internal server error page.
func Route500(c *gin.Context) {
	c.HTML(http.StatusInternalServerError, "500.html", nil)
}

// Health returns a simple health check response.
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]string{"status": "OK"})
}

// ServiceWorker serves the service worker JavaScript file.
func ServiceWorker(c *gin.Context) {
	c.Header("Cache-Control", "public, max-age=31536000, immutable")
	c.Header("ETag", faviconETag)
	c.File("frontend/service-worker.js")
}

// OfflinePage renders the offline fallback page.
func OfflinePage(c *gin.Context) {
	c.HTML(http.StatusOK, "offline_page.html", nil)
}
