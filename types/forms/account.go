// Package forms defines request form structs used for binding and validation.
package forms

// AccountEnable represents the form for enabling an account via hash.
type AccountEnable struct {
	Hash string `form:"hash" json:"hash" query:"hash" binding:"required,uuid4"`
}
