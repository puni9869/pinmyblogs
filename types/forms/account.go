package forms

type AccountEnableForm struct {
	Hash string `form:"hash" json:"hash" query:"hash" binding:"required,uuid4"`
}
