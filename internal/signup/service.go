package signup

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/puni9869/pinmyblogs/models"
	"github.com/puni9869/pinmyblogs/server/middlewares"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var ErrDuplicateEmail = errors.New("Email already exist")

type Service interface {
	Register(ctx *gin.Context, user models.User) error
	CheckPassword()
	CheckEmail()
	IsActive()
	Verify()
	IsVerified()
}

type signupClient struct {
	db  *gorm.DB
	log *logrus.Logger
}

func (s *signupClient) Register(c *gin.Context, user models.User) error {
	ctx := middlewares.GetContext(c)
	err := s.db.Create(&user).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		s.log.WithError(err).Error("failed to create user")
		ctx["Email_HasError"] = true
		ctx["Email_Err"] = ErrDuplicateEmail.Error()
		ctx["Password_HasError"] = false
		ctx["ConfirmPassword_HasError"] = false
		return ErrDuplicateEmail
	}
	if err != nil {
		s.log.WithError(err).Error("failed to create user")
		return err
	}
	s.log.Infoln("user is created successfully")
	return nil
}

func (s *signupClient) CheckPassword() {
	s.log.Infoln("CheckPassword")
}

func (s *signupClient) CheckEmail() {
	s.log.Infoln("CheckEmail")
}

func (s *signupClient) IsActive() {
	s.log.Infoln("IsActive")
}

func (s *signupClient) Verify() {
	s.log.Infoln("Verify")
}

func (s *signupClient) IsVerified() {
}

func NewSignupService(db *gorm.DB, logger *logrus.Logger) Service {
	return &signupClient{
		db:  db,
		log: logger,
	}
}
