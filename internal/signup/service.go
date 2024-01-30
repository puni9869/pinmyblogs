package signup

import (
	"gitea.com/go-chi/binding"
	"github.com/puni9869/pinmyblogs/types/forms"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
)

// SetForm set the form object
func SetForm(dataStore map[string]any, obj any) {
	dataStore["__form"] = obj
}

// GetForm returns the validate form information
func GetForm(dataStore map[string]any) any {
	return dataStore["__form"]
}

type SignupService interface {
	Register()
	ValidateForm(r *http.Request) binding.Errors
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

func (c *signupClient) ValidateForm(r *http.Request) binding.Errors {
	// create a new form obj for every request but not use obj directly
	theObj := new(forms.SignUpForm)
	var data = new(map[string]any)
	binding.Bind(r, theObj)
	c.log.Infof("%T", theObj)
	SetForm(data, theObj)
	//formbinding.AssignForm(theObj, data)
	//var errs binding.Errors
	return nil
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
