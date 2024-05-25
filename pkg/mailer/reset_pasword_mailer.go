package mailer

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/puni9869/pinmyblogs/models"
	"github.com/puni9869/pinmyblogs/pkg/config"
	"github.com/puni9869/pinmyblogs/pkg/logger"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type ResetPasswordMailer struct {
	MailerAPI
	log  *logrus.Logger
	user models.User
}

func (u *ResetPasswordMailer) getPasswordResetLink() string {
	h, _ := uuid.NewUUID()
	return fmt.Sprintf("https://pinmyblogs.com/reset/password/%s", h.String())
}

func (u *ResetPasswordMailer) Send() {
	if config.GetEnv() == config.LocalEnv {
		u.log.WithFields(map[string]any{
			"user": u.user.Email,
			"env":  config.GetEnv(),
		}).Info("reset password mail has been sent")
		return
	}

	tmpl := fmt.Sprintf(`<!DOCTYPE>
	<html>
	<body>
       Hi there, 
	   <br/>
	   To reset your Pinmyblogs account password, please click the line below <br/> %s

       <br/>
       Regards,
		<br/>
       pinmyblogs and team
	</body>
	</html>
	`, u.getPasswordResetLink())

	msg := gomail.NewMessage()
	msg.SetHeader("From", config.C.Mailer.EmailId)
	msg.SetHeader("To", u.user.Email)
	msg.SetHeader("Bcc", config.C.Mailer.BccEmailId)
	msg.SetHeader("Subject", "Reset Password")
	msg.SetBody("text/html", tmpl)

	m := gomail.NewDialer(
		config.C.Mailer.SmtpHost,
		config.C.Mailer.SmtpPort,
		config.C.Mailer.Username,
		config.C.Mailer.Password,
	)

	if err := m.DialAndSend(msg); err != nil {
		u.log.WithError(err).Error("error sending email for reset password")
		return
	}
	u.log.WithFields(map[string]any{
		"user": u.user.Email,
		"id":   u.user.ID,
		"env":  config.GetEnv(),
	}).Info("reset password email has been send")
}
func NewResetPasswordMailer(user models.User) MailerAPI {
	return &ResetPasswordMailer{
		log:  logger.NewLogger(),
		user: user,
	}
}
