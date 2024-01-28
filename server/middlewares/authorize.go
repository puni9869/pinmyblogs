package middlewares

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

const userkey = "user"

// AuthRequired is a simple middleware to check the session.
func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	fmt.Println("AuthRequired")
	if user == nil {
		// Abort the request with the appropriate error code
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}
	// Continue down the chain to handler etc
	c.Next()
}
