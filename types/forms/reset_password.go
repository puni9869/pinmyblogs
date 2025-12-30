package forms

type ResetPasswordForm struct {
	Hash            string `form:"hash" json:"hash" binding:"required,hash"`
	Password        string `form:"password" json:"password" binding:"required"`
	ConfirmPassword string `form:"confirm_password" json:"confirm_password" binding:"required"`
}
