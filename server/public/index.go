package public

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/puni9869/pinmyblogs"
	"github.com/puni9869/pinmyblogs/models"
	"github.com/puni9869/pinmyblogs/pkg/database"
	"github.com/puni9869/pinmyblogs/pkg/logger"
	"github.com/puni9869/pinmyblogs/server/middlewares"
	"github.com/puni9869/pinmyblogs/types/forms"
)

var (
	favicon     []byte
	faviconETag string
)

func StartGet(c *gin.Context) {
	c.HTML(http.StatusOK, "start.tmpl", nil)
}
func StartPost(c *gin.Context) {
	log := logger.NewLogger()
	form := middlewares.GetForm(c).(*forms.JoinWaitList)
	email := form.Email
	ctx := middlewares.GetContext(c)
	if ctx["Email_HasError"] == true {
		ctx["HasError"] = true
		log.WithFields(map[string]any{"email": email}).
			WithError(errors.New(ctx["Email_Error"].(string))).
			Error("bad email from join wait list page.")

		c.HTML(http.StatusBadRequest, "start.tmpl", ctx)
		return
	}

	j := models.JoinWaitList{Email: email, App: "pinmyblogs"}
	err := database.Db().Save(&j).Error
	if err != nil {
		log.WithFields(map[string]any{
			"email": j.Email, "app": j.App,
		}).WithError(err).Error("failed to save user into wait list")
		c.HTML(http.StatusBadRequest, "start.tmpl", ctx)
		return
	}
	log.Infof("%s has joined the wait list", email)
	c.HTML(http.StatusAccepted, "join_wait_list_pinmyblogs.tmpl", ctx)
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

func Route5xx(c *gin.Context) {
	c.HTML(http.StatusInternalServerError, "500.tmpl", nil)
}
