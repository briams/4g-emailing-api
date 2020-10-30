package validators

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// CreateBody payload for update method
type CreateBody struct {
	ModelID   string `json:"modelId"`
	Mjml      string `json:"mjml"`
	Variables string `json:"variables"`
	InsUserID uint   `json:"insUserId"`
}

// Validate method for CreateBody
func (b *CreateBody) Validate() error {
	return validation.ValidateStruct(b,
		validation.Field(&b.ModelID, validation.Required),
		validation.Field(&b.Mjml, validation.Required),
		validation.Field(&b.Variables, validation.Required),
		validation.Field(&b.InsUserID, validation.Required),
	)
}
