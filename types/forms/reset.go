package forms

type ResetForm struct {
	Email string `form:"email" json:"email" binding:"required,email"`
}
