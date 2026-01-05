package forms

type ResetPassword struct {
	Hash            string `form:"hash" json:"hash" binding:"required"`
	Password        string `form:"password" json:"password" binding:"required"`
	ConfirmPassword string `form:"confirm_password" json:"confirm_password" binding:"required"`
}
