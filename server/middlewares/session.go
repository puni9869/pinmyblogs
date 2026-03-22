package middlewares

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/gorm"
	"github.com/gin-gonic/gin"
)

// Session returns a middleware that manages user sessions using the given store.
func Session(sessionStore gorm.Store) gin.HandlerFunc {
	return sessions.Sessions("session", sessionStore)
}
