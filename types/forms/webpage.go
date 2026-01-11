package forms

type WeblinkRequest struct {
	Url string `form:"url" json:"url" binding:"required"`
	Tag string `form:"tag" json:"tag" binding:"required"`
}

type BulkAction struct {
	IDs    []string `form:"ids[]" json:"ids" binding:"required"`
	Action string   `form:"action" json:"action" binding:"required"`
}
