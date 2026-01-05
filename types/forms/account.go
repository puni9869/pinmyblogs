package forms

type AccountEnable struct {
	Hash string `form:"hash" json:"hash" query:"hash" binding:"required,uuid4"`
}
