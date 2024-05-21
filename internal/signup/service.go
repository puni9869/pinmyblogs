package signup

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/puni9869/pinmyblogs/models"
	"github.com/puni9869/pinmyblogs/pkg/mailer"
	"github.com/puni9869/pinmyblogs/server/middlewares"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var ErrDuplicateEmail = errors.New("Email already exist")

type Service interface {
	Register(ctx *gin.Context, user models.User) error
	Verify()
	Notify(user models.User)
}

type signupClient struct {
	db  *gorm.DB
	log *logrus.Logger
}

func (s *signupClient) Register(c *gin.Context, user models.User) error {
	ctx := middlewares.GetContext(c)
	user.IsActive = true
	user.IsEmailVerified = true
	user.ActivatedAt = time.Now()
	err := s.db.Create(&user).Error
	if err != nil {
		s.log.WithError(err).Error("failed to create user")
		ctx["Email_HasError"] = true
		ctx["HasError"] = true
		ctx["Email_Error"] = ErrDuplicateEmail.Error()
		ctx["Password_HasError"] = false
		ctx["ConfirmPassword_HasError"] = false
		return ErrDuplicateEmail
	}
	s.log.Infoln("user is created successfully")
	return nil
}

func (s *signupClient) Verify() {
	s.log.Infoln("user is verified")
}

func (s *signupClient) Notify(user models.User) {
	userMailer := mailer.NewUserRegisterMailer(user)
	go userMailer.Send()
}

func NewSignupService(db *gorm.DB, logger *logrus.Logger) Service {
	return &signupClient{
		db:  db,
		log: logger,
	}
}
