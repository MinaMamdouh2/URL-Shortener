// Package auth provides authentication and authorization support.
// Authentication: You are who you say you are.
// Authorization:  You have permission to do what you are requesting to do.
package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// ErrForbidden is returned when a auth issue is identified.
// The idiom for those error variables is you start with "Err", and it is exported because we want to access it.
// The decision to export so we can do something like this: "errors.Is(err, auth.ErrForbidden) to return HTTP 403."
var ErrForbidden = errors.New("attempted action is not allowed")

// Claims represents the authorization claims transmitted via a JWT.
type Claims struct {
	// The subject in the claims will be the uuid for the user
	jwt.RegisteredClaims
	Roles []string `json:"roles"`
}

// A RuleFn returns nil if authorization passes, or an error if it fails.
type ruleFn func(claims Claims) error

// ruleFns maps rule names (as passed into Authorize) to their Go implementations.
var ruleFns = map[string]ruleFn{
	"admin_only": adminOnly,
}

// KeyLookup declares a method set of behavior for looking up private and public keys for JWT use.
// The return could be a PEM encoded string or a JWS based key.
// This interface we are gonna use it to ask the application layer how we are gonna get the public and private keys.
type KeyLookup interface {
	PrivateKey(kid string) (key string, err error)
	PublicKey(kid string) (key string, err error)
}

// Config represents information required to initialize auth.
// In the config we are asking for a logger, we can do that because we are in the business layer.
// We have the implementation of the KeyLookup interface and the Issuer name.
// The "Config" carries dependency injection that will come from the App layer
type Config struct {
	Log       *zap.SugaredLogger
	KeyLookup KeyLookup
	Issuer    string
}

// Auth is used to authenticate clients.
// It can generate a token for a set of user claims and recreate the claims by parsing the token.
type Auth struct {
	log       *zap.SugaredLogger
	keyLookup KeyLookup
	method    jwt.SigningMethod
	parser    *jwt.Parser
	issuer    string
	// this basic caching is for saving the keys so we don't have to do the network call to get the key every time
	// In reality, you need to implement something to periodically fetch the keys back again.
	mu    sync.RWMutex
	cache map[string]string
}

// New creates an Auth to support authentication/authorization.
func New(cfg Config) (*Auth, error) {
	a := Auth{
		keyLookup: cfg.KeyLookup,
		method:    jwt.GetSigningMethod(jwt.SigningMethodRS256.Name),
		parser:    jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Name})),
		issuer:    cfg.Issuer,
	}

	return &a, nil
}

// GenerateToken generates a signed JWT token string representing the user Claims.
func (a *Auth) GenerateToken(kid string, claims Claims) (string, error) {
	// Creating a token using claims with kid
	token := jwt.NewWithClaims(a.method, claims)
	token.Header["kid"] = kid

	privateKeyPEM, err := a.keyLookup.PrivateKey(kid)
	if err != nil {
		return "", fmt.Errorf("private key: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKeyPEM))
	if err != nil {
		return "", fmt.Errorf("parsing private pem: %w", err)
	}

	str, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("signing token: %w", err)
	}

	return str, nil
}

// Authenticate processes the token to validate the sender's token is valid.
func (a *Auth) Authenticate(ctx context.Context, bearerToken string) (Claims, error) {
	parts := strings.Split(bearerToken, " ")
	// 1. Split “Bearer <jwt>”
	if len(parts) != 2 || parts[0] != "Bearer" {
		return Claims{}, errors.New("expected authorization header format: Bearer <token>")
	}

	// 2. ParseUnverified to extract claims without checking signature
	var claims Claims
	token, _, err := a.parser.ParseUnverified(parts[1], &claims)
	if err != nil {
		return Claims{}, fmt.Errorf("error parsing token: %w", err)
	}

	// 3. Extract string `kid` from header
	kidRaw, exists := token.Header["kid"]
	if !exists {
		return Claims{}, fmt.Errorf("kid missing from header: %w", err)
	}

	// 4. Doing a type assertion to make sure it is a string
	// Go asks here "Is the underlying value inside this interface{} actually a string?"
	kid, ok := kidRaw.(string)
	if !ok {
		return Claims{}, fmt.Errorf("kid malformed: %w", err)
	}

	// 5. Gets public key
	pem, err := a.publicKeyLookup(kid)
	if err != nil {
		return Claims{}, fmt.Errorf("failed to fetch public key: %w", err)
	}

	keyFunc := func(t *jwt.Token) (interface{}, error) {
		// Parse the PEM into an *rsa.PublicKey
		return jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
	}

	// 6. Validate token
	_, err = a.parser.ParseWithClaims(parts[1], &claims, keyFunc)

	if err != nil {
		return Claims{}, fmt.Errorf("authentication failed : %w", err)
	}

	// 7. Check if the issuer is a valid issuer
	if claims.Issuer != a.issuer {
		return Claims{}, fmt.Errorf("authentication failed: wrong issuer")
	}

	// Check the database for this user to verify they are still enabled.
	if err := a.isUserEnabled(ctx, claims); err != nil {
		return Claims{}, fmt.Errorf("user not enabled : %w", err)
	}

	return claims, nil
}

// Authorize attempts to authorize the user with the provided input roles, if
// none of the input roles are within the user's claims, we return an error
// otherwise the user is authorized.
func (a *Auth) Authorize(claims Claims, userID uuid.UUID, rule string) error {

	fn, ok := ruleFns[rule]
	if !ok {
		return fmt.Errorf("authorization rule not found")
	}

	// Executing the rule function
	if err := fn(claims); err != nil {
		return err
	}

	return nil
}

// =============================================================================
// This barer for unexported functions

// publicKeyLookup performs a lookup for the public pem for the specified kid.
func (a *Auth) publicKeyLookup(kid string) (string, error) {
	// Bill wanted to make this more simpler, so he made a literal function and executed it right away.
	pem, err := func() (string, error) {
		a.mu.RLock()
		defer a.mu.RUnlock()
		// This is a cache lookup for the pem key
		pem, exists := a.cache[kid]
		if !exists {
			return "", errors.New("not found")
		}
		return pem, nil
	}()

	if err == nil {
		return pem, nil
	}
	// if it wasn't there we will make the network call "assuming that is on another service"
	pem, err = a.keyLookup.PublicKey(kid)

	if err != nil {
		return "", fmt.Errorf("fetching public key: %w", err)
	}
	// Store in the cache
	a.mu.Lock()
	defer a.mu.Unlock()
	a.cache[kid] = pem
	return pem, nil
}

func adminOnly(claims Claims) error {
	for _, v := range claims.Roles {
		if v == "admin_only" {
			return nil
		}
	}
	return ErrForbidden
}

// isUserEnabled hits the database and checks the user is not disabled. If the
// no database connection was provided, this check is skipped.
func (a *Auth) isUserEnabled(ctx context.Context, claims Claims) error {

	return nil
}
