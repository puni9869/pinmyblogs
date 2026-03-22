package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestCors_SetsHeaders(t *testing.T) {
	r := gin.New()
	r.Use(Cors())
	r.GET("/api", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api", nil)
	req.Header.Set("Origin", "https://example.com")
	r.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "https://example.com" {
		t.Errorf("Allow-Origin = %q, want https://example.com", w.Header().Get("Access-Control-Allow-Origin"))
	}
	if w.Header().Get("Access-Control-Allow-Credentials") != "true" {
		t.Error("Allow-Credentials should be true")
	}
	if w.Header().Get("Access-Control-Allow-Methods") == "" {
		t.Error("Allow-Methods should be set")
	}
}

func TestCors_OptionsReturns200(t *testing.T) {
	r := gin.New()
	r.Use(Cors())
	r.OPTIONS("/api", func(c *gin.Context) {
		c.String(http.StatusOK, "should not reach")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodOptions, "/api", nil)
	req.Header.Set("Origin", "https://example.com")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("OPTIONS status = %d, want 200", w.Code)
	}
}

func TestCors_NonOptionsPassesThrough(t *testing.T) {
	r := gin.New()
	r.Use(Cors())
	r.POST("/api", func(c *gin.Context) {
		c.String(http.StatusCreated, "created")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("POST status = %d, want 201", w.Code)
	}
}
