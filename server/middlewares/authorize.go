package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const Userkey = "user"

func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(Userkey)
	if user == nil {
		// Redirect to the login page if not authenticated
		contentType := c.GetHeader("Content-Type")
		if strings.EqualFold(contentType, "application/json") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Status": "NOT_OK", "Errors": "NOT_AUTHENTICATED"})
		}
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		c.Abort()
	}
	// Continue down the chain to the handler, etc.
	c.Handler()
}
