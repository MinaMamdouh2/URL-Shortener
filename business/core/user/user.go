// Bill uses singular names in the business and the foundation layer. In app layer he tends to use plural names more.
package user

import (
	"context"
	"fmt"
	"net/mail"
	"time"

	"github.com/MinaMamdouh2/URL-Shortener/business/data/order"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// Storer interface declares the behavior this package needs to persists and retrieve data.
// Bill always focused on this idea of small interface that provide the sort of this generic functionality,
// your reader, your writer, he is always focused on the idea of precision.
// We shouldn't have interfaces that are incredibly large because a lot of times there isn't a single function
// that is necessarily going to need all of that behavior.
// Also, interfaces should describe verbs "their behavior", they are not things but these ideas always break down
// when we start wanting to abstract physical things "socket, file system, DB",
// what is happening is we want to abstract things now not just behavior,
// you will end up with larger interfaces and it is ok.
// Any activity, that we need to perform against the datastore, we are going to define that behavior here.
type Storer interface {
	// Create: Given the user type, we expect that is stored in the DB.
	Create(ctx context.Context, usr User) error
	// Update: Given the user type, we expect that is updated in the DB.
	Update(ctx context.Context, usr User) error
	// Update: Given the user type, we expect that is deleted in the DB.
	Delete(ctx context.Context, usr User) error
	// Query, it is a generic query function this eliminates the possibility of lots of query functions.
	// also "Query", should be the only query to return a collection.
	//  All other functions that starts with "Query", returns a singular. It also provides paging.
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, pageNumber int, rowsPerPage int) ([]User, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, userID uuid.UUID) (User, error)
	// This is an exception, for returning a collection because we had this special case
	// that we need to return specific users by IDs and we wanted to optimize for that.
	// This doesn't provide paging and it is very dangerous API.
	QueryByIDs(ctx context.Context, userID []uuid.UUID) ([]User, error)
	QueryByEmail(ctx context.Context, email mail.Address) (User, error)
}

// =============================================================================
// Core manages the set of APIs for user access.
// In every single core business package, we define a type named core
// which will represent the API for that business package.
type Core struct {
	storer Storer
	log    *zap.SugaredLogger
}

// NewCore constructs a core for user api access.
func NewCore(log *zap.SugaredLogger, storer Storer) *Core {
	return &Core{
		storer: storer,
		log:    log,
	}
}

// Create adds a new user to the system.
// We are using pointer semantics because core doesn't represent data, it represents an API.
// We return User as a value because we represent that user data.
// Also User is value semantics because it represents pure data.
func (c *Core) Create(ctx context.Context, nu NewUser) (User, error) {
	// bcrypt to generate hash
	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)

	if err != nil {
		return User{}, fmt.Errorf("generatefrompassword: %w", err)
	}

	// This will be the createdAt time
	now := time.Now()

	usr := User{
		ID:           uuid.New(),
		Name:         nu.Name,
		Email:        nu.Email,
		Roles:        nu.Roles,
		PasswordHash: hash,
		Enabled:      true,
		DateCreated:  now,
		DateUpdated:  now,
	}

	if err := c.storer.Create(ctx, usr); err != nil {
		return User{}, fmt.Errorf("create user: %w", err)
	}

	return User{}, nil
}
