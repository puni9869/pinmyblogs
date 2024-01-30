package forms

import (
	"gitea.com/go-chi/binding"
	"github.com/puni9869/pinmyblogs/pkg/formbinding"
	"net/http"
)

type SignUpForm struct {
	Email           string `form:"email" json:"email" binding:"required,email"`
	Password        string `form:"password" json:"password" binding:"required"`
	ConfirmPassword string `form:"confirm_password" json:"confirm_password" binding:"required"`
}

// Validate validates fields
func (f *SignUpForm) Validate(req *http.Request, errs binding.Errors) binding.Errors {
	return formbinding.Validate(errs, f)
}
