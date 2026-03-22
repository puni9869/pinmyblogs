// Package blankeditor provides handlers for the blank-blankeditor SPA page.
package blankeditor

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/puni9869/pinmyblogs"
)

// EditorGet serves the blank-editor SPA page by reading the
// Vite-built index.html directly from the embedded filesystem.
func EditorGet(c *gin.Context) {
	data, err := pinmyblogs.Files.ReadFile("frontend/blankeditor/index.html")
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", data)
}
