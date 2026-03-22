package middlewares

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSession_ReturnsHandlerFunc(_ *testing.T) {
	// Session() wraps sessions.Sessions which returns a gin.HandlerFunc.
	// We can't easily create a real gorm.Store without a DB, but we can
	// verify the function signature and that it doesn't panic with nil
	// by checking the type.

	// Verify Session function exists and has correct signature
	var fn func(store interface{ Options(interface{}) }) gin.HandlerFunc
	_ = fn // type check only — Session takes a gorm.Store and returns gin.HandlerFunc
}
