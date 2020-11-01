package validators

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// ActiveLogBody payload for update method
type ActiveLogBody struct {
	SetUserID uint   `json:"setUserId"`
	Reason    string `json:"reason"`
}

// Validate method for ActiveLogBody
func (b *ActiveLogBody) Validate() error {
	return validation.ValidateStruct(
		b,
		validation.Field(&b.SetUserID, validation.Required),
		validation.Field(&b.Reason, validation.Required),
	)
}
