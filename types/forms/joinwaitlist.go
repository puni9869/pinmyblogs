package forms

// JoinWaitList represents the form for joining the wait list.
type JoinWaitList struct {
	Email string `form:"email" json:"email" binding:"required,email"`
}
