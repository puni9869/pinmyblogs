package mailer

import (
	"github.com/puni9869/pinmyblogs/models"
	"github.com/puni9869/pinmyblogs/pkg/logger"
	"github.com/sirupsen/logrus"
)

type UserRegister struct {
	MailerAPI
	log  *logrus.Logger
	user models.User
}

func (u *UserRegister) Send() {
	u.log.WithFields(map[string]any{
		"user": u.user.Email,
		"id":   u.user.ID,
	}).Info("user registration confirmation mail sent")

}

func NewUserRegisterMailer(user models.User) MailerAPI {
	return &UserRegister{
		log:  logger.NewLogger(),
		user: user,
	}
}
