package forms

// Reset represents the form for requesting a password reset.
type Reset struct {
	Email string `form:"email" json:"email" binding:"required,email"`
}
