package public

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestFetchPublicUserByUsername_ValidUser(t *testing.T) {
	user, err := fetchPublicUserByUsername("alex41")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Username != "alex41" {
		t.Errorf("Username = %q, want %q", user.Username, "alex41")
	}
	if user.JoinedAt != "January 2025" {
		t.Errorf("JoinedAt = %q, want %q", user.JoinedAt, "January 2025")
	}
	if len(user.PublicBlogs) != 2 {
		t.Errorf("PublicBlogs count = %d, want 2", len(user.PublicBlogs))
	}
	if user.Bio == "" {
		t.Error("Bio should not be empty")
	}
}

func TestFetchPublicUserByUsername_NotFound(t *testing.T) {
	_, err := fetchPublicUserByUsername("nonexistent")
	if err == nil {
		t.Error("expected error for unknown user")
	}
}

func TestFetchPublicUserByUsername_EmptyUsername(t *testing.T) {
	_, err := fetchPublicUserByUsername("")
	if err == nil {
		t.Error("expected error for empty username")
	}
}

func TestFetchPublicUserByUsername_TooLong(t *testing.T) {
	long := "abcdefghijklmnopqrstuvwxyz12345" // 31 chars
	_, err := fetchPublicUserByUsername(long)
	if err == nil {
		t.Error("expected error for username > 30 chars")
	}
}

func TestFetchPublicUserByUsername_InvalidChars(t *testing.T) {
	invalids := []string{"user name", "user@name", "User", "UPPER", "hello!"}
	for _, u := range invalids {
		t.Run(u, func(t *testing.T) {
			_, err := fetchPublicUserByUsername(u)
			if err == nil {
				t.Errorf("expected error for invalid username %q", u)
			}
		})
	}
}

func TestFetchPublicUserByUsername_ValidChars(t *testing.T) {
	// Only "alex41" exists in the stub, but we can test that
	// valid-format usernames that don't exist return "user not found"
	valids := []string{"valid-user", "user_name", "abc123"}
	for _, u := range valids {
		t.Run(u, func(t *testing.T) {
			_, err := fetchPublicUserByUsername(u)
			if err == nil {
				t.Errorf("expected 'user not found' for %q (not in stub)", u)
			}
			if err.Error() != "user not found" {
				t.Errorf("error = %q, want %q", err.Error(), "user not found")
			}
		})
	}
}

func TestFetchPublicUserByUsername_TrimsAndLowercases(t *testing.T) {
	// "ALEX41" with spaces should be lowercased and trimmed to "alex41"
	user, err := fetchPublicUserByUsername("  ALEX41  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Username != "alex41" {
		t.Errorf("Username = %q, want %q", user.Username, "alex41")
	}
}

func TestFetchPublicUserByUsername_ExactlyAt30Chars(t *testing.T) {
	// 30 chars is the boundary — should pass validation but not be found
	name := "aaaaabbbbbcccccdddddeeeeefffff" // exactly 30
	_, err := fetchPublicUserByUsername(name)
	if err == nil {
		t.Error("expected error (user not found), got nil")
	}
	if err.Error() != "user not found" {
		t.Errorf("error = %q, want %q", err.Error(), "user not found")
	}
}

func TestUserPublicProfilePage_NonProfilePath(t *testing.T) {
	r := gin.New()
	r.NoRoute(UserPublicProfilePage)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/some-random-path", nil)
	r.ServeHTTP(w, req)

	// Should try to render 404.html — since no template is loaded, Gin returns 500
	// but the handler logic is exercised
	if w.Code == http.StatusOK {
		t.Error("non-profile path should not return 200")
	}
}

func TestUserPublicProfilePage_InvalidUsername(t *testing.T) {
	r := gin.New()
	r.NoRoute(UserPublicProfilePage)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/@INVALID!", nil)
	r.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		t.Error("invalid username should not return 200")
	}
}

func TestUserRepoFindByUsername_Stub(t *testing.T) {
	user, err := userRepoFindByUsername("alex41")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID != 1 {
		t.Errorf("ID = %d, want 1", user.ID)
	}
	if user.DisplayName != "alex41" {
		t.Errorf("DisplayName = %q, want %q", user.DisplayName, "alex41")
	}

	_, err = userRepoFindByUsername("nobody")
	if err == nil {
		t.Error("expected error for unknown username")
	}
}

func TestBlogRepoFindPublicByUserID_Stub(t *testing.T) {
	blogs, err := blogRepoFindPublicByUserID()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(blogs) != 2 {
		t.Errorf("blogs count = %d, want 2", len(blogs))
	}
	if blogs[0].Domain != "blog.example.com" {
		t.Errorf("first blog domain = %q, want %q", blogs[0].Domain, "blog.example.com")
	}
}
