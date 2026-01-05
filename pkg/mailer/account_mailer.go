package mailer

import (
	"fmt"
	"github.com/puni9869/pinmyblogs/models"
	"github.com/puni9869/pinmyblogs/pkg/config"
	"github.com/puni9869/pinmyblogs/pkg/logger"
	"github.com/sirupsen/logrus"
	m "gopkg.in/gomail.v2"
)

type Account struct {
	MailerAPI
	log    *logrus.Logger
	user   models.User
	action string
}

func (a *Account) disableAccount() (*m.Dialer, *m.Message) {
	tmpl := fmt.Sprintf(`<!DOCTYPE>
	<html>
	<body>
		<p> This is to confirm that your PinMyBlogs account associated with the email <strong>%s</strong> has been successfully disabled. </p>
       
		<p> Your data remains secure, and your account will stay inactive until you choose to log in again. </p>

		<p> If this action wasn’t intended or you have any questions, feel free to reach out to us. </p> 
		
		<p> Regards,<br /> <strong>PinMyBlogs Team</strong> </p>

	</body>
	</html>
	`, a.user.Email)

	msg := m.NewMessage()
	msg.SetHeader("From", config.C.Mailer.EmailId)
	msg.SetHeader("To", a.user.Email)
	msg.SetHeader("Bcc", config.C.Mailer.BccEmailId)
	msg.SetHeader("Subject", "Your pinmyblogs.com Account Has Been Disabled")
	msg.SetBody("text/html", tmpl)

	dialer := m.NewDialer(
		config.C.Mailer.SmtpHost,
		config.C.Mailer.SmtpPort,
		config.C.Mailer.Username,
		config.C.Mailer.Password,
	)

	return dialer, msg
}

func (a *Account) enableAccount() (*m.Dialer, *m.Message) {
	host := map[string]string{"local": "http://localhost", "prod": "https://pinmyblogs.com"}[config.GetEnv()]
	tmpl := fmt.Sprintf(`<!DOCTYPE>
			<html>
		  <body>
			<p>Hello there,</p>
			<p>
			  You previously disabled your <strong>PinMyBlogs</strong> account associated with the email
			  <strong>%s</strong>.
			</p>
			<p>If you’d like to re-enable your account, simply click the button below:</p>
			<p style="margin: 24px 0">
			  <a
				href="%s/enable-my-account/%s"
				style="
				  display: inline-block;
				  background-color: #4f46e5;
				  color: #ffffff;
				  padding: 12px 20px;
				  text-decoration: none;
				  border-radius: 8px;
				  font-weight: 600;
				"
			  >
				Enable My Account
			  </a>
			</p>
			<p>For your security, this link is unique and will expire after one time use.</p>
			<p>If you didn’t request this or believe this email was sent by mistake, you can safely ignore it.</p>
			<p>
			  Regards,<br />
			  <strong>PinMyBlogs Team</strong>
			</p>
		  </body>
		</html>
	`, a.user.Email, host, a.user.AccountEnableHash)

	msg := m.NewMessage()
	msg.SetHeader("From", config.C.Mailer.EmailId)
	msg.SetHeader("To", a.user.Email)
	msg.SetHeader("Bcc", config.C.Mailer.BccEmailId)
	msg.SetHeader("Subject", "Enable Your pinmyblogs.com Account")
	msg.SetBody("text/html", tmpl)

	dialer := m.NewDialer(
		config.C.Mailer.SmtpHost,
		config.C.Mailer.SmtpPort,
		config.C.Mailer.Username,
		config.C.Mailer.Password,
	)

	return dialer, msg
}

func (a *Account) Send() {
	f := map[string]any{
		"user": a.user.Email,
		"id":   a.user.ID,
		"env":  config.GetEnv(),
	}
	var msg *m.Message
	var m *m.Dialer
	switch a.action {
	case "disable":
		m, msg = a.disableAccount()
	case "enable":
		m, msg = a.enableAccount()
	}

	if config.GetEnv() == config.LocalEnv {
		a.log.WithFields(f).Infof("user account %s confirmation mail sent. hash %s", a.action, a.user.AccountEnableHash)
		return
	}

	if err := m.DialAndSend(msg); err != nil {
		a.log.WithFields(f).WithError(err).Errorf("error sending email for account %s confirmation", a.action)
		return
	}

	a.log.WithFields(f).Infof("user account %s confirmation mail sent. hash %s", a.action, a.user.AccountEnableHash)
}

func NewAccountService(user models.User, action string) MailerAPI {
	return &Account{
		log:    logger.NewLogger(),
		user:   user,
		action: action,
	}
}
