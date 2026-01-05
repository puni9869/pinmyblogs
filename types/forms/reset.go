package forms

type Reset struct {
	Email string `form:"email" json:"email" binding:"required,email"`
}
