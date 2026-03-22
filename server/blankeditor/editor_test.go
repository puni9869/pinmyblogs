package blankeditor

import (
	"io/fs"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"path"
	"regexp"
	"strings"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/puni9869/pinmyblogs"
	"github.com/puni9869/pinmyblogs/server/middlewares"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// setupRouter creates a Gin engine with session + auth middleware and the blankeditor route.
// The authenticate parameter controls whether a session user is injected before AuthRequired.
func setupRouter(authenticate bool) *gin.Engine {
	r := gin.New()
	store := cookie.NewStore([]byte("test-secret-key-for-sessions"))
	r.Use(sessions.Sessions("session", store))

	if authenticate {
		// Inject a session user so AuthRequired passes.
		r.Use(func(c *gin.Context) {
			s := sessions.Default(c)
			s.Set(middlewares.Userkey, "testuser")
			_ = s.Save()
			c.Next()
		})
	}

	r.Use(middlewares.AuthRequired)
	r.GET("/blankeditor", EditorGet)
	return r
}

// randomUserID generates a random user identifier string for property testing.
func randomUserID(rng *rand.Rand) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	length := rng.Intn(20) + 1
	b := make([]byte, length)
	for i := range b {
		b[i] = chars[rng.Intn(len(chars))]
	}
	return string(b)
}

// TestProperty_AuthenticatedEditorRouteReturnsEditorHTML validates Property 1:
// For any authenticated HTTP request to the blankeditor route, the Go server responds
// with HTTP 200 and a body containing a complete standalone HTML document that
// does NOT contain pinmyblogs layout markers.
//
// **Validates: Requirements 1.1, 3.1**
// Feature: blank-blankeditor-integration, Property 1: Authenticated blankeditor route returns blankeditor HTML
func TestProperty_AuthenticatedEditorRouteReturnsEditorHTML(t *testing.T) {
	const iterations = 100

	rng := rand.New(rand.NewSource(42))

	for i := 0; i < iterations; i++ {
		userID := randomUserID(rng)

		r := gin.New()
		store := cookie.NewStore([]byte("test-secret-key-for-sessions"))
		r.Use(sessions.Sessions("session", store))
		r.Use(func(c *gin.Context) {
			s := sessions.Default(c)
			s.Set(middlewares.Userkey, userID)
			_ = s.Save()
			c.Next()
		})
		r.Use(middlewares.AuthRequired)
		r.GET("/blankeditor", EditorGet)

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/blankeditor", nil)
		r.ServeHTTP(w, req)

		// Property: HTTP status must be 200
		if w.Code != http.StatusOK {
			t.Fatalf("iteration %d (user=%q): status = %d, want 200", i, userID, w.Code)
		}

		body := w.Body.String()
		bodyLower := strings.ToLower(body)

		// Property: body contains required HTML document tags
		requiredTags := []string{"<!doctype html>", "<html", "<head", "<body"}
		for _, tag := range requiredTags {
			if !strings.Contains(bodyLower, tag) {
				t.Fatalf("iteration %d (user=%q): body missing %q", i, userID, tag)
			}
		}

		// Property: body does NOT contain pinmyblogs layout markers
		forbiddenMarkers := []string{"side-navbar"}
		for _, marker := range forbiddenMarkers {
			if strings.Contains(bodyLower, marker) {
				t.Fatalf("iteration %d (user=%q): body contains pinmyblogs layout marker %q", i, userID, marker)
			}
		}

		// Property: Content-Type is text/html
		ct := w.Header().Get("Content-Type")
		if !strings.Contains(ct, "text/html") {
			t.Fatalf("iteration %d (user=%q): Content-Type = %q, want text/html", i, userID, ct)
		}
	}
}

// mimeForExt returns the expected Content-Type for a given file extension.
// Returns empty string for extensions where Go's HTTP file server does not
// register a specific MIME type (e.g., .webmanifest), in which case the test
// skips Content-Type validation but still verifies HTTP 200.
func mimeForExt(ext string) string {
	switch strings.ToLower(ext) {
	case ".js":
		return "text/javascript"
	case ".css":
		return "text/css"
	case ".html":
		return "text/html"
	case ".png":
		return "image/png"
	case ".svg":
		return "image/svg+xml"
	case ".ico":
		return "image/x-icon"
	case ".json":
		return "application/json"
	case ".txt":
		return "text/plain"
	case ".xml":
		return "text/xml"
	default:
		return ""
	}
}

// setupRouterWithStatics creates a Gin engine with session + auth middleware,
// the blankeditor handler route, and the static file server for blankeditor assets.
func setupRouterWithStatics(userID string) *gin.Engine {
	r := gin.New()
	store := cookie.NewStore([]byte("test-secret-key-for-sessions"))
	r.Use(sessions.Sessions("session", store))
	r.Use(func(c *gin.Context) {
		s := sessions.Default(c)
		s.Set(middlewares.Userkey, userID)
		_ = s.Save()
		c.Next()
	})
	r.Use(middlewares.AuthRequired)
	r.GET("/blankeditor", EditorGet)

	// Serve blankeditor static assets at /blankeditor/statics/ from the embedded FS.
	editorFS, err := fs.Sub(pinmyblogs.Files, "frontend/blankeditor")
	if err != nil {
		panic("failed to create blankeditor sub-filesystem: " + err.Error())
	}
	r.StaticFS("/blankeditor/statics", http.FS(editorFS))

	return r
}

// TestProperty_AllEditorAssetReferencesAreServable validates Property 2:
// For any asset URL (script src, link href, img src) referenced in the blankeditor's
// HTML page, the Go server serves that asset with HTTP 200 and a Content-Type
// header matching the file extension.
//
// **Validates: Requirements 2.3, 2.4, 3.2, 5.1, 5.2, 9.2, 9.3**
// Feature: blank-blankeditor-integration, Property 2: All blankeditor HTML asset references are servable with correct MIME types
func TestProperty_AllEditorAssetReferencesAreServable(t *testing.T) {
	const iterations = 100

	// Regex to extract src and href attribute values that start with /blankeditor/statics/
	attrRe := regexp.MustCompile(`(?:src|href)="(/blankeditor/statics/[^"]+)"`)

	rng := rand.New(rand.NewSource(42))

	for i := 0; i < iterations; i++ {
		userID := randomUserID(rng)
		r := setupRouterWithStatics(userID)

		// Step 1: GET /blankeditor to obtain the HTML page.
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/blankeditor", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("iteration %d (user=%q): GET /blankeditor status = %d, want 200", i, userID, w.Code)
		}

		body := w.Body.String()

		// Step 2: Extract all asset URLs from src and href attributes.
		matches := attrRe.FindAllStringSubmatch(body, -1)
		if len(matches) == 0 {
			t.Fatalf("iteration %d (user=%q): no /blankeditor/statics/ asset references found in HTML", i, userID)
		}

		// Step 3: Request each asset and verify HTTP 200 + correct Content-Type.
		for _, m := range matches {
			assetURL := m[1]
			ext := path.Ext(assetURL)
			expectedMIME := mimeForExt(ext)

			aw := httptest.NewRecorder()
			areq := httptest.NewRequest(http.MethodGet, assetURL, nil)
			r.ServeHTTP(aw, areq)

			if aw.Code != http.StatusOK {
				t.Fatalf("iteration %d (user=%q): GET %s status = %d, want 200",
					i, userID, assetURL, aw.Code)
			}

			ct := aw.Header().Get("Content-Type")
			if expectedMIME != "" && !strings.Contains(ct, expectedMIME) {
				t.Fatalf("iteration %d (user=%q): GET %s Content-Type = %q, want to contain %q",
					i, userID, assetURL, ct, expectedMIME)
			}
		}
	}
}
