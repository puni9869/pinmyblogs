package mailer

import (
	"fmt"

	"github.com/puni9869/pinmyblogs/models"
	"github.com/puni9869/pinmyblogs/pkg/config"
	"github.com/puni9869/pinmyblogs/pkg/logger"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type UserRegister struct {
	MailerAPI
	log  *logrus.Logger
	user models.User
}

func (u *UserRegister) Send() {
	if config.GetEnv() == config.LocalEnv {
		u.log.WithFields(map[string]any{
			"user": u.user.Email,
			"id":   u.user.ID,
			"env":  config.GetEnv(),
		}).Info("user registration confirmation mail sent")
		return
	}
	tmpl := fmt.Sprintf(`<!DOCTYPE>
	<html>
	<body>
       Hi there, welcome to pinmyblogs.
	   We have successfully registered the user with email: %s<br/>
       <br/>
       Regards,
		<br/>
       pinmyblogs and team
	</body>
	</html>
	`, u.user.Email)

	msg := gomail.NewMessage()
	msg.SetHeader("From", config.C.Mailer.EmailId)
	msg.SetHeader("To", u.user.Email)
	msg.SetHeader("Bcc", config.C.Mailer.BccEmailId)
	msg.SetHeader("Subject", "Welcome to pinmyblogs.com")
	msg.SetBody("text/html", tmpl)

	m := gomail.NewDialer(
		config.C.Mailer.SmtpHost,
		config.C.Mailer.SmtpPort,
		config.C.Mailer.Username,
		config.C.Mailer.Password,
	)

	if err := m.DialAndSend(msg); err != nil {
		u.log.WithError(err).Error("error sending email")
		return
	}
	u.log.WithFields(map[string]any{
		"user": u.user.Email,
		"id":   u.user.ID,
		"env":  config.GetEnv(),
	}).Info("user registration confirmation mail sent")
}

func NewUserRegisterMailer(user models.User) MailerAPI {
	return &UserRegister{
		log:  logger.NewLogger(),
		user: user,
	}
}
