package validators

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// CreateBody payload for update method
type CreateBody struct {
	Name      string `json:"name"`
	InsUserID uint   `json:"insUserId"`
}

// Validate method for CreateBody
func (b *CreateBody) Validate() error {
	return validation.ValidateStruct(b,
		validation.Field(&b.Name, validation.Required,
			validation.Match(regexp.MustCompile("^[a-zA-Z0-9 ]*$")).Error("only must have numbers, letters and spaces"),
		),
		validation.Field(&b.InsUserID, validation.Required),
	)
}
