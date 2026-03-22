package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestCacheMiddleware_StaticFiles(t *testing.T) {
	extensions := []string{".js", ".css", ".png", ".jpg", ".svg", ".ico", ".woff2"}
	for _, ext := range extensions {
		t.Run(ext, func(t *testing.T) {
			r := gin.New()
			r.Use(CacheMiddleware())
			r.GET("/statics/file"+ext, func(c *gin.Context) {
				c.String(http.StatusOK, "ok")
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/statics/file"+ext, nil)
			r.ServeHTTP(w, req)

			cc := w.Header().Get("Cache-Control")
			if cc != "public, max-age=86400, must-revalidate" {
				t.Errorf("Cache-Control for %s = %q, want public cache header", ext, cc)
			}
		})
	}
}

func TestCacheMiddleware_NonStaticFiles(t *testing.T) {
	paths := []string{"/home", "/login", "/api/data"}
	for _, p := range paths {
		t.Run(p, func(t *testing.T) {
			r := gin.New()
			r.Use(CacheMiddleware())
			r.GET(p, func(c *gin.Context) {
				c.String(http.StatusOK, "ok")
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, p, nil)
			r.ServeHTTP(w, req)

			cc := w.Header().Get("Cache-Control")
			if cc != "no-cache" {
				t.Errorf("Cache-Control for %s = %q, want no-cache", p, cc)
			}
		})
	}
}

func TestCSP_AppliesHeader(t *testing.T) {
	r := gin.New()
	r.Use(CSP([]string{"/exempt"}))
	r.GET("/protected", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	r.ServeHTTP(w, req)

	csp := w.Header().Get("Content-Security-Policy")
	if csp == "" {
		t.Error("CSP header should be set on non-exempted route")
	}
	if len(csp) < 20 {
		t.Errorf("CSP header seems too short: %q", csp)
	}
}

func TestCSP_ExemptedRoute(t *testing.T) {
	r := gin.New()
	r.Use(CSP([]string{"/exempt"}))
	r.GET("/exempt", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/exempt", nil)
	r.ServeHTTP(w, req)

	csp := w.Header().Get("Content-Security-Policy")
	if csp != "" {
		t.Errorf("CSP header should NOT be set on exempted route, got %q", csp)
	}
}

func TestSecurityHeaders_Applied(t *testing.T) {
	r := gin.New()
	r.Use(SecurityHeaders([]string{}))
	r.GET("/page", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/page", nil)
	r.ServeHTTP(w, req)

	expected := map[string]string{
		"X-Content-Type-Options":       "nosniff",
		"X-Frame-Options":              "DENY",
		"Referrer-Policy":              "strict-origin-when-cross-origin",
		"X-XSS-Protection":             "0",
		"Cross-Origin-Opener-Policy":   "same-origin",
		"Cross-Origin-Resource-Policy": "same-origin",
	}
	for header, want := range expected {
		got := w.Header().Get(header)
		if got != want {
			t.Errorf("%s = %q, want %q", header, got, want)
		}
	}

	hsts := w.Header().Get("Strict-Transport-Security")
	if hsts == "" {
		t.Error("Strict-Transport-Security should be set")
	}

	pp := w.Header().Get("Permissions-Policy")
	if pp == "" {
		t.Error("Permissions-Policy should be set")
	}
}

func TestSecurityHeaders_ExemptedRoute(t *testing.T) {
	r := gin.New()
	r.Use(SecurityHeaders([]string{"/skip"}))
	r.GET("/skip", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/skip", nil)
	r.ServeHTTP(w, req)

	if w.Header().Get("X-Frame-Options") != "" {
		t.Error("Security headers should NOT be set on exempted route")
	}
}
