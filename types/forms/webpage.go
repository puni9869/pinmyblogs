package forms

// WeblinkRequest represents the form for adding a new weblink.
type WeblinkRequest struct {
	Url string `form:"url" json:"url" binding:"required"`
	Tag string `form:"tag" json:"tag" binding:"required"`
}

// BulkAction represents the form for performing bulk operations on URLs.
type BulkAction struct {
	IDs    []string `form:"ids[]" json:"ids" binding:"required"`
	Action string   `form:"action" json:"action" binding:"required"`
}
