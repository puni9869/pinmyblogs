package public

import (
	"errors"
	"net/http"

	"github.com/puni9869/pinmyblogs/types/forms"

	"github.com/gin-gonic/gin"
	"github.com/puni9869/pinmyblogs/models"
	"github.com/puni9869/pinmyblogs/pkg/database"
	"github.com/puni9869/pinmyblogs/pkg/logger"
	"github.com/puni9869/pinmyblogs/server/middlewares"
)

func JoinWaitListPost(c *gin.Context) {
	log := logger.NewLogger()
	form := middlewares.GetForm(c).(forms.JoinWaitList)
	email := form.Email
	ctx := middlewares.GetContext(c)
	if ctx["Email_HasError"] == true {
		ctx["HasError"] = true
		log.WithFields(map[string]any{"email": email}).
			WithError(errors.New(ctx["Email_Error"].(string))).
			Error("bad email from join wait list page.")

		c.HTML(http.StatusBadRequest, "join_wait_list_pinmyblogs.html", ctx)
		return
	}

	j := models.JoinWaitList{Email: email, App: "pinmyblogs"}
	err := database.Db().Save(&j).Error
	if err != nil {
		log.WithFields(map[string]any{
			"email": j.Email, "app": j.App,
		}).WithError(err).Error("failed to save user into wait list")
		c.HTML(http.StatusBadRequest, "join_wait_list_pinmyblogs.html", ctx)
		return
	}
	log.Infof("%s has joined the wait list", email)
	c.HTML(http.StatusAccepted, "join_wait_list_pinmyblogs.html", ctx)
}

func JoinWaitListGet(c *gin.Context) {
	c.HTML(http.StatusOK, "join_wait_list_pinmyblogs.html", nil)
}
