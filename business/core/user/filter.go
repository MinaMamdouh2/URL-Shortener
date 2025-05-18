package user

import (
	"fmt"
	"net/mail"
	"time"

	"github.com/MinaMamdouh2/URL-Shortener/foundation/validate"
	"github.com/google/uuid"
)

// QueryFilter holds the available fields a query can be filtered on.
// The QueryFilter type, has been defined in the business layer,
// you could also see we are using pointer semantics which means we are leveraging pointer semantics
// which means that concept of null, you may or may not want to filter on these given fields
// and we also added a validate tag, for use for out validate APIs.
// This allows the App layer when it calls our validate APIs it gets applied to this.
type QueryFilter struct {
	ID               *uuid.UUID    `validate:"omitempty"` // skip the validation if the pointer is nil
	Name             *string       `validate:"omitempty,min=3"`
	Email            *mail.Address `validate:"omitempty"`
	StartCreatedDate *time.Time    `validate:"omitempty"`
	EndCreatedDate   *time.Time    `validate:"omitempty"`
}

// Validate checks the data in the model is considered clean.
// Calls our foundation validation.Check which reads the "validate:" tags.
// If any filter is set but invalid, you get back a "FieldErrors" error describing which fields failed.
func (qf *QueryFilter) Validate() error {
	// validate knows how to read the tags and do the validation against the data model
	if err := validate.Check(qf); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	return nil
}

// We provide an entire API for each of those fields and we use the concept of with.
// So, you construct a query filter and you say filter with the userID, name ...,
// that gives us that data is set properly.

// WithUserID sets the ID field of the QueryFilter value.
func (qf *QueryFilter) WithUserID(userID uuid.UUID) {
	qf.ID = &userID
}

// WithName sets the Name field of the QueryFilter value.
func (qf *QueryFilter) WithName(name string) {
	qf.Name = &name
}

// WithEmail sets the Email field of the QueryFilter value.
func (qf *QueryFilter) WithEmail(email mail.Address) {
	qf.Email = &email
}

// WithStartDateCreated sets the DateCreated field of the QueryFilter value.
func (qf *QueryFilter) WithStartDateCreated(startDate time.Time) {
	d := startDate.UTC()
	qf.StartCreatedDate = &d
}

// WithEndCreatedDate sets the DateCreated field of the QueryFilter value.
func (qf *QueryFilter) WithEndCreatedDate(endDate time.Time) {
	d := endDate.UTC()
	qf.EndCreatedDate = &d
}
