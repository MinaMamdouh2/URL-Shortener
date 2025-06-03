// This is were we define the data models,
// or those data types that are going to represent data coming in and data coming out.
package user

import (
	"net/mail"
	"time"

	"github.com/google/uuid"
)

// user represents information about an individual user.
// This is the user model it represents the data, that we need to maintain from the business perspective on a user
// and what it also does and it needs to do as part of the idea of forcing trust,
// is use types that have their own forum of parsing were business logic behind it,
// parsing is a big key to enforce trust, so the UUID package is going to be our ID.
// One of the things you have to decide early on is how unique IDs are generated for different entities
// that we are creating.
// You want the DB to generate it or you, we want here the business package to do it.
// We have from the UUID a type that we can enforce validation on it, but for the scalar types like "string"
// we don't have that sort of framework to do this validation, for example we are saying that "Name" could be anything because it has a string,
// but we trust the app layer to have done the validation that the Name is valid.
// We introduce this type Role, we want a user to be locked in a static set of roles where roles would be a string.
// If we made it a string, we are not forcing the app layer nor we are giving it any tooling to make sure that they give
// me a role that is going to be well defined in the package, we are going to define a type.
// Also if any business model had a set of tags for marshaling that is a smell, that enforces us to know that these
// models are not what get marshaled and unmarshaled.

type User struct {
	ID           uuid.UUID    `gorm:"column:id;type:uuid;primaryKey"`
	Name         string       `gorm:"column:name"`
	Email        mail.Address `gorm:"column:email"`
	Roles        []Role       `gorm:"column:roles;type:text[]"` // TODO: Look into pq.StringArray to replace []Role
	PasswordHash []byte       `gorm:"column:password_has"`
	Enabled      bool         `gorm:"column:enabled"`
	DateCreated  time.Time    `gorm:"column:date_created"`
	DateUpdated  time.Time    `gorm:"column:date_updated"`
}

// NewUser contains information needed to create a new user.
// Bill likes to have a model, for creation and he does that because if he wants to use the user model,
// it is not obvious what field you need to provide and which that you don't.
// He doesn't want the caller or App layer to construct a UIID,
// we are going to do that, we don't need the caller to create the dates we are going to do that.
type NewUser struct {
	Name            string
	Email           mail.Address
	Roles           []Role
	Password        string
	PasswordConfirm string
}

// UpdateUser contains information needed to update a user.
// We are using pointer semantics, updates a very complicated thing in a data system.
// It is nice from the usability standpoint to be able to ask the user to specify what it is that they want to update,
// instead of potentially asking them to get a complete record of something, make the changes and then do it.
// By using pointer semantics, we say if you just want to update the email, leave everything nil
// and we are gonna update that.
// We are using pointer semantics to represent the concept of null.
type UpdateUser struct {
	Name            *string
	Email           *mail.Address
	Roles           []Role
	Password        *string
	PasswordConfirm *string
	Enabled         *bool
}
