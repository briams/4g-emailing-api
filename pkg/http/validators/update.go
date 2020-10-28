package validators

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// UpdateBody payload for update method
type UpdateBody struct {
	Name      string `json:"name"`
	SetUserID uint   `json:"setUserId"`
}

// Validate method for UpdateBody
func (b *UpdateBody) Validate() error {
	return validation.ValidateStruct(b,
		validation.Field(&b.Name, validation.Required, validation.Match(regexp.MustCompile("^[a-zA-Z0-9 ]*$"))),
		validation.Field(&b.SetUserID, validation.Required),
	)
}
