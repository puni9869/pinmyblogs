package signup

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SignupService interface {
	Register()
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

func (c *signupClient) Register() {
	c.log.Infoln("Register")
}

func (c *signupClient) CheckPassword() {
	c.log.Infoln("CheckPassword")
}

func (c *signupClient) CheckEmail() {
	c.log.Infoln("CheckEmail")
}

func (c *signupClient) IsActive() {
	c.log.Infoln("IsActive")
}

func (c *signupClient) Verify() {
	c.log.Infoln("Verify")
}

func (c *signupClient) IsVerified() {
}

func NewSignupService(db *gorm.DB, logger *logrus.Logger) SignupService {
	return &signupClient{
		db:  db,
		log: logger,
	}
}
