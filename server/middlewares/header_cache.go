package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
	"strings"
	"time"
)

// CacheMiddleware adds cache headers to static file responses
func CacheMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ext := strings.ToLower(path.Ext(c.Request.URL.Path))
		switch ext {
		case ".js", ".css", ".png", ".jpg", ".jpeg", ".svg", ".gif", ".ico", ".woff2", ".map":
			c.Writer.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			c.Writer.Header().Set("Expires", time.Now().AddDate(1, 0, 0).Format(http.TimeFormat))
		default:
			c.Writer.Header().Set("Cache-Control", "no-cache")
		}
		c.Next()
	}
}
