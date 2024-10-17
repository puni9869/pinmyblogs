package middlewares

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

const Userkey = "user"

func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(Userkey)
	if user == nil {
		// Redirect to the login page if not authenticated
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		c.Abort()
	}
	// Continue down the chain to the handler, etc.
	c.Handler()
}
