package middlewares

import (
	"path"
	"strings"

	"github.com/gin-gonic/gin"
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

// CSP adds a Content-Security-Policy (CSP) header to every HTTP response.
//
// This middleware helps prevent Cross-Site Scripting (XSS) and related
// injection attacks by restricting where resources can be loaded from.
//
// Policy applied:
//   - default-src 'self'        → Allow resources only from the same origin
//   - script-src 'self'         → Allow scripts only from the same origin
//   - style-src 'self' 'unsafe-inline'
//     → Allow styles from the same origin and
//     inline styles (required for Tailwind / inline CSS)
//
// Note:
//   - 'unsafe-inline' is intentionally allowed for styles but NOT for scripts
//   - For higher security, consider replacing inline styles with hashes or nonces
//   - This middleware should be registered early in the Gin middleware chain
func CSP() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Security-Policy",
			"default-src 'self'; "+
				"script-src 'self'; "+
				"style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; "+
				"font-src 'self' https://fonts.gstatic.com; "+
				"img-src 'self' data: https://www.google.com; "+
				"connect-src 'self'; "+
				"object-src 'none'; "+
				"frame-ancestors 'none'",
		)
		c.Next()
	}
}

// SecurityHeaders adds common HTTP security headers to every response.
//
// These headers help protect the application against a range of attacks
// including MIME sniffing, clickjacking, and information leakage.
//
// Headers applied:
//
//   - X-Content-Type-Options: nosniff
//     Prevents browsers from MIME-sniffing responses away from the declared content-type.
//
//   - X-Frame-Options: DENY
//     Prevents the site from being embedded in iframes (clickjacking protection).
//
//   - Referrer-Policy: strict-origin-when-cross-origin
//     Limits the amount of referrer information sent with cross-origin requests.
//
//   - X-XSS-Protection: 0
//     Explicitly disables legacy browser XSS filters, which are deprecated and
//     can introduce security issues in modern browsers.
//
// Note:
//   - This middleware should be applied globally.
//   - Works best when combined with a strong Content Security Policy (CSP).
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("X-XSS-Protection", "0")

		// HTTPS only
		c.Header(
			"Strict-Transport-Security",
			"max-age=63072000; includeSubDomains; preload",
		)

		// Modern security headers
		c.Header(
			"Permissions-Policy",
			"geolocation=(), microphone=(), camera=(), payment=()",
		)
		c.Header("Cross-Origin-Opener-Policy", "same-origin")
		c.Header("Cross-Origin-Resource-Policy", "same-origin")

		c.Next()
	}
}
