package middlewares

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

const userkey = "user"

func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		// Redirect to the login page if not authenticated
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		c.Abort()
	}
	// Continue down the chain to the handler, etc.
	c.Handler()
}
