package forms

// Prefs represents the form for updating user preferences.
type Prefs struct {
	Action string `form:"action" json:"action" binding:"required,min=2,max=20"`
	Value  string `form:"value" json:"value" binding:"required,min=2,max=20"`
}
