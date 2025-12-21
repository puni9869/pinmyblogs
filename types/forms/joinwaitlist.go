package forms

type JoinWaitList struct {
	Email string `form:"email" json:"email" binding:"required,email"`
}
