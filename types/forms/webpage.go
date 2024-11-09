package forms

type WeblinkRequest struct {
	Url string `form:"url" json:"url" binding:"required"`
	Tag string `form:"tag" json:"tag" binding:"required"`
}
