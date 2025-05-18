// This is a good pattern to use when you have to enforce a validation at the app layer
package user

import "fmt"

// Set of possible roles for a user.
// We construct some exported variables, that represent the set of roles that are valid
// because what we are doing in there is constructing a value of type role.
// Since this is an Exported variable, we can change it outside this package,
// yes but we trust the team and we can catch it.
// It can't be a constant since we need it as type Role
var (
	RoleAdmin = Role{"ADMIN"}
	RoleUser  = Role{"USER"}
)

// Set of known roles.
// Here we get some on enforcement, we define a map of roles based on the name and the value of the role itself.
var roles = map[string]Role{
	RoleAdmin.name: RoleAdmin,
	RoleUser.name:  RoleUser,
}

// Role represents a role in the system.
// We are going to define a type Role, and it has an unexported field called "name", the key it is unexported.
// The fact that is unexported means that the app layer can not construct a role outside of zero value
// that is not proper, because they can't stuff any string they want in there.
// Zero value is the only hole in this approach.
type Role struct {
	name string
}

// ParseRole parses the string value and returns a role if one exists.
// We define the parse function, this forces the app layer to use the
// parse function to construct a role and it does the validation.
// Now if they pass any string to this, if it doesn't map to our map it is not valid,
// if it does we get a valid role
func ParseRole(value string) (Role, error) {
	role, exists := roles[value]
	if !exists {
		return Role{}, fmt.Errorf("invalid role %q", value)
	}

	return role, nil
}

// MustParseRole parses the string value and returns a role if one exists.
// If an error occurs the function panics.
// For testing only, we add a Must, the Must functions only being used in tests to streamline the construction of a User,
// of some role that's why it panics.
func MustParseRole(value string) Role {
	role, err := ParseRole(value)
	if err != nil {
		panic(err)
	}

	return role
}

// Name returns the name of the role.
// You have to be able to get that name.
// Notice we are using value semantics, why because this is pure data.
// No pointers here.
func (r Role) Name() string {
	return r.name
}

// Here for "UnmarshalText" & "MarshalText", we allow the marshaler to marshal this
// into whatever the different formats are.
// If we implement the MarshalText and UnmarshalText functions,
// your JSON and another Marshalers will see this implementation and use it.
// You don't have to write a marshaler for everyone, only write Text and you get all the others
// for free and we need this because when we marshal a role back into a JSON.

// Go’s encoding packages look for MarshalText/UnmarshalText when(de)serializing arbitrary types.
// Implementing these lets Role integrate smoothly.

// UnmarshalText implement the unmarshal interface for JSON conversions.
func (r *Role) UnmarshalText(data []byte) error {
	// Unmarshal calls parse, that means also at the app layer when a value of some type of this field,
	// if there is an invalid string it will fail
	role, err := ParseRole(string(data))
	if err != nil {
		return err
	}

	r.name = role.name
	return nil
}

// MarshalText implement the marshal interface for JSON conversions.
func (r Role) MarshalText() ([]byte, error) {
	return []byte(r.name), nil
}

// Equal provides support for the go-cmp package and testing.
// The Equal method is there, because Bill uses a package from google called "cmp" for testing
// and this allows this package to compare between two roles.
func (r Role) Equal(r2 Role) bool {
	return r.name == r2.name
}

// TODO: Add an Is method to validate if the a value of type Role is a zero value, but we don't have the need for that
// because Bill doesn't believe in code that you write to protect yourself from yourself.
