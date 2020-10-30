package validators

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// UpdateBody payload for update method
type UpdateBody struct {
	Mjml      string `json:"mjml"`
	Variables string `json:"variables"`
	SetUserID uint   `json:"setUserId"`
}

// Validate method for UpdateBody
func (b *UpdateBody) Validate() error {
	return validation.ValidateStruct(b,
		validation.Field(&b.Mjml, validation.Required),
		validation.Field(&b.Variables, validation.Required),
		validation.Field(&b.SetUserID, validation.Required),
	)
}
