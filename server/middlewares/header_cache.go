package middlewares

import (
	"github.com/gin-gonic/gin"
	"path"
	"strings"
)

// CacheMiddleware adds cache headers to static file responses
func CacheMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ext := strings.ToLower(path.Ext(c.Request.URL.Path))
		switch ext {
		case ".js", ".css", ".png", ".jpg", ".jpeg", ".svg", ".gif", ".ico", ".woff2", ".map":
			c.Writer.Header().Set("Cache-Control", "public, max-age=86400, must-revalidate")
		default:
			c.Writer.Header().Set("Cache-Control", "no-cache")
		}
		c.Next()
	}
}
