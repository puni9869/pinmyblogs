package forms

type Prefs struct {
	Action string `form:"action" json:"action" binding:"required,min=2,max=20"`
	Value  string `form:"value" json:"value" binding:"required,min=2,max=20"`
}
