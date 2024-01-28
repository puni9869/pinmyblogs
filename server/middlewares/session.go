package middlewares

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/gorm"
	"github.com/gin-gonic/gin"
)

func Session(sessionStore gorm.Store) gin.HandlerFunc {
	return sessions.Sessions("session", sessionStore)
}
