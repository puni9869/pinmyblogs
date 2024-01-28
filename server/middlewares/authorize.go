package middlewares

import (
	"fmt"
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
		return
	}
	// Continue down the chain to the handler, etc.
	fmt.Println(c.HandlerNames())
	c.Handler()
	return
}
