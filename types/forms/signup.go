package forms

type SignUpForm struct {
	Email           string `form:"email" json:"email" binding:"required,email"`
	Password        string `form:"password" json:"password" binding:"required"`
	ConfirmPassword string `form:"confirm_password" json:"confirm_password" binding:"required"`
}
