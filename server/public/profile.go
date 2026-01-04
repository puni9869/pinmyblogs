package public

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// UserPublicProfilePage profile page
func UserPublicProfilePage(c *gin.Context) {
	path := c.Request.URL.Path

	// Only handle /@username
	if !strings.HasPrefix(path, "/@") {
		c.HTML(http.StatusNotFound, "404.tmpl", nil)
		return
	}

	username := strings.TrimPrefix(path, "/@")
	username = strings.TrimSpace(username)

	user, err := fetchPublicUserByUsername(username)
	if err != nil {
		c.HTML(http.StatusNotFound, "404.tmpl", gin.H{
			"Username": username,
		})
		return
	}

	c.HTML(http.StatusOK, "profile_page.tmpl", gin.H{
		"User":  user,
		"Blogs": user.PublicBlogs,
	})
}

//
// ================= DOMAIN =================
//

type User struct {
	ID          int64
	Username    string
	DisplayName string
	Bio         string
	AvatarURL   string
	CreatedAt   time.Time
}

type PublicUser struct {
	Username    string
	DisplayName string
	Bio         string
	AvatarURL   string
	JoinedAt    string
	PublicBlogs []PublicBlog
}

type PublicBlog struct {
	Title       string
	URL         string
	Description string
	Domain      string
}

//
// ================= FETCH LOGIC =================
//

func fetchPublicUserByUsername(username string) (*PublicUser, error) {
	username = strings.ToLower(strings.TrimSpace(username))

	if username == "" || len(username) > 30 {
		return nil, errors.New("invalid username")
	}

	// allow a-z 0-9 - _
	for _, r := range username {
		if !(r >= 'a' && r <= 'z') &&
			!(r >= '0' && r <= '9') &&
			r != '-' && r != '_' {
			return nil, errors.New("invalid username")
		}
	}

	user, err := userRepoFindByUsername(username)
	if err != nil {
		return nil, err
	}

	blogs, err := blogRepoFindPublicByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	return &PublicUser{
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Bio:         user.Bio,
		AvatarURL:   user.AvatarURL,
		JoinedAt:    user.CreatedAt.Format("January 2006"),
		PublicBlogs: blogs,
	}, nil
}

//
// ================= REPOSITORIES (STUBS) =================
//

func userRepoFindByUsername(username string) (*User, error) {
	if username != "alex41" {
		return nil, errors.New("user not found")
	}

	return &User{
		ID:          1,
		Username:    "alex41",
		DisplayName: "alex41",
		Bio:         "Quietly collecting thoughtful blogs on software and systems.",
		AvatarURL:   "https://via.placeholder.com/120",
		CreatedAt:   time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC),
	}, nil
}

func blogRepoFindPublicByUserID(userID int64) ([]PublicBlog, error) {
	return []PublicBlog{
		{
			Title:       "PostgreSQL Indexing Explained",
			URL:         "https://blog.example.com/postgres-indexing",
			Description: "A clear breakdown of how indexes work.",
			Domain:      "blog.example.com",
		},
		{
			Title:       "Designing Distributed Systems",
			URL:         "https://medium.com/systems/designing-distributed-systems",
			Description: "Trade-offs and real-world lessons.",
			Domain:      "medium.com",
		},
	}, nil
}
